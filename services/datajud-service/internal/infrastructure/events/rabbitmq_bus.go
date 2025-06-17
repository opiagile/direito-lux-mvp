package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
	"github.com/direito-lux/datajud-service/internal/infrastructure/logging"
	"github.com/direito-lux/datajud-service/internal/infrastructure/metrics"
	"github.com/direito-lux/datajud-service/internal/infrastructure/tracing"
)

// RabbitMQEventBus implementação do event bus usando RabbitMQ
type RabbitMQEventBus struct {
	config     *config.Config
	logger     *zap.Logger
	metrics    *metrics.Metrics
	connection *amqp.Connection
	channel    *amqp.Channel
	handlers   map[string][]EventHandler
	serializer EventSerializer
	mu         sync.RWMutex
	done       chan bool
}

// NewRabbitMQEventBus cria novo event bus com RabbitMQ
func NewRabbitMQEventBus(cfg *config.Config, logger *zap.Logger, metrics *metrics.Metrics) (*RabbitMQEventBus, error) {
	bus := &RabbitMQEventBus{
		config:     cfg,
		logger:     logger,
		metrics:    metrics,
		handlers:   make(map[string][]EventHandler),
		serializer: NewJSONEventSerializer(),
		done:       make(chan bool),
	}

	if err := bus.connect(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com RabbitMQ: %w", err)
	}

	return bus, nil
}

// connect estabelece conexão com RabbitMQ
func (bus *RabbitMQEventBus) connect() error {
	var err error
	
	// Conectar
	bus.connection, err = amqp.Dial(bus.config.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("erro ao conectar: %w", err)
	}

	// Criar canal
	bus.channel, err = bus.connection.Channel()
	if err != nil {
		return fmt.Errorf("erro ao criar canal: %w", err)
	}

	// Configurar QoS
	if err := bus.channel.Qos(bus.config.RabbitMQ.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("erro ao configurar QoS: %w", err)
	}

	// Declarar exchange
	if err := bus.channel.ExchangeDeclare(
		bus.config.RabbitMQ.Exchange,
		"topic",
		bus.config.RabbitMQ.Durable,
		bus.config.RabbitMQ.AutoDelete,
		bus.config.RabbitMQ.Exclusive,
		bus.config.RabbitMQ.NoWait,
		nil,
	); err != nil {
		return fmt.Errorf("erro ao declarar exchange: %w", err)
	}

	// Declarar queue
	_, err = bus.channel.QueueDeclare(
		bus.config.RabbitMQ.Queue,
		bus.config.RabbitMQ.Durable,
		bus.config.RabbitMQ.AutoDelete,
		bus.config.RabbitMQ.Exclusive,
		bus.config.RabbitMQ.NoWait,
		amqp.Table{
			"x-dead-letter-exchange":    "direito_lux.dlx",
			"x-dead-letter-routing-key": bus.config.RabbitMQ.Queue + ".dlq",
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar queue: %w", err)
	}

	// Binding
	if err := bus.channel.QueueBind(
		bus.config.RabbitMQ.Queue,
		bus.config.RabbitMQ.RoutingKey + ".*",
		bus.config.RabbitMQ.Exchange,
		bus.config.RabbitMQ.NoWait,
		nil,
	); err != nil {
		return fmt.Errorf("erro ao fazer binding: %w", err)
	}

	bus.logger.Info("Event bus conectado ao RabbitMQ",
		zap.String("exchange", bus.config.RabbitMQ.Exchange),
		zap.String("queue", bus.config.RabbitMQ.Queue),
		zap.String("routing_key", bus.config.RabbitMQ.RoutingKey),
	)

	// Monitorar conexão
	go bus.monitorConnection()

	return nil
}

// monitorConnection monitora a conexão e reconecta se necessário
func (bus *RabbitMQEventBus) monitorConnection() {
	for {
		select {
		case <-bus.done:
			return
		case err := <-bus.connection.NotifyClose(make(chan *amqp.Error)):
			if err != nil {
				bus.logger.Error("Conexão RabbitMQ perdida", zap.Error(err))
				
				// Tentar reconectar
				for {
					bus.logger.Info("Tentando reconectar ao RabbitMQ...")
					if err := bus.connect(); err != nil {
						bus.logger.Error("Erro ao reconectar", zap.Error(err))
						time.Sleep(bus.config.RabbitMQ.ReconnectDelay)
						continue
					}
					
					bus.logger.Info("Reconectado ao RabbitMQ com sucesso")
					break
				}
			}
		}
	}
}

// Publish publica evento no RabbitMQ
func (bus *RabbitMQEventBus) Publish(ctx context.Context, event DomainEvent) error {
	span, ctx := tracing.TraceMessage(ctx, "publish", bus.config.RabbitMQ.Exchange, bus.config.RabbitMQ.RoutingKey)
	defer span.Finish()

	// Serializar evento
	data, err := bus.serializer.Serialize(event)
	if err != nil {
		tracing.SetSpanError(span, err)
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	// Preparar headers
	headers := amqp.Table{
		"event_type":   event.GetEventType(),
		"event_id":     event.GetEventID(),
		"aggregate_id": event.GetAggregateID(),
		"tenant_id":    event.GetTenantID(),
		"occurred_at":  event.GetOccurredAt().Format(time.RFC3339),
		"version":      event.GetVersion(),
	}

	// Adicionar trace headers
	if traceHeaders := make(map[string]string); tracing.InjectHTTPHeaders(ctx, traceHeaders) == nil {
		for k, v := range traceHeaders {
			headers["trace_"+k] = v
		}
	}

	// Routing key específico por tipo de evento
	routingKey := fmt.Sprintf("%s.%s", bus.config.RabbitMQ.RoutingKey, event.GetEventType())

	// Publicar
	err = bus.channel.Publish(
		bus.config.RabbitMQ.Exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    event.GetEventID(),
			Timestamp:    event.GetOccurredAt(),
			Headers:      headers,
			Body:         data,
		},
	)

	// Registrar métricas
	success := err == nil
	if bus.metrics != nil {
		bus.metrics.RecordMessageSent(
			bus.config.RabbitMQ.Exchange,
			routingKey,
			event.GetTenantID(),
			success,
		)
	}

	if err != nil {
		tracing.SetSpanError(span, err)
		logging.LogError(ctx, bus.logger, "Erro ao publicar evento", err,
			zap.String("event_type", event.GetEventType()),
			zap.String("event_id", event.GetEventID()),
			zap.String("routing_key", routingKey),
		)
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}

	logging.LogInfo(ctx, bus.logger, "Evento publicado",
		zap.String("event_type", event.GetEventType()),
		zap.String("event_id", event.GetEventID()),
		zap.String("routing_key", routingKey),
	)

	return nil
}

// Subscribe subscreve handler para tipo de evento
func (bus *RabbitMQEventBus) Subscribe(eventType string, handler EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = make([]EventHandler, 0)
	}
	
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	
	bus.logger.Info("Handler subscrito",
		zap.String("event_type", eventType),
	)
	
	return nil
}

// Start inicia o consumo de mensagens
func (bus *RabbitMQEventBus) Start(ctx context.Context) error {
	// Consumir mensagens
	deliveries, err := bus.channel.Consume(
		bus.config.RabbitMQ.Queue,
		"datajud-service", // consumer tag
		false,              // auto-ack
		bus.config.RabbitMQ.Exclusive,
		false, // no-local
		bus.config.RabbitMQ.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("erro ao iniciar consumo: %w", err)
	}

	bus.logger.Info("Event bus iniciado, consumindo mensagens",
		zap.String("queue", bus.config.RabbitMQ.Queue),
	)

	// Processar mensagens
	go bus.processMessages(ctx, deliveries)

	return nil
}

// processMessages processa mensagens recebidas
func (bus *RabbitMQEventBus) processMessages(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for delivery := range deliveries {
		go bus.handleMessage(ctx, delivery)
	}
}

// handleMessage processa uma mensagem individual
func (bus *RabbitMQEventBus) handleMessage(ctx context.Context, delivery amqp.Delivery) {
	start := time.Now()
	
	// Extrair informações da mensagem
	eventType, ok := delivery.Headers["event_type"].(string)
	if !ok {
		bus.logger.Error("Tipo de evento não encontrado nos headers")
		delivery.Nack(false, false) // Enviar para DLQ
		return
	}

	eventID, _ := delivery.Headers["event_id"].(string)
	tenantID, _ := delivery.Headers["tenant_id"].(string)

	// Criar contexto com informações do evento
	if tenantID != "" {
		ctx = logging.WithTenantID(ctx, tenantID)
	}

	// Extrair trace context
	traceHeaders := make(map[string]string)
	for k, v := range delivery.Headers {
		if key, ok := k.(string); ok && len(key) > 6 && key[:6] == "trace_" {
			if value, ok := v.(string); ok {
				traceHeaders[key[6:]] = value
			}
		}
	}
	
	if len(traceHeaders) > 0 {
		if spanCtx, err := tracing.ExtractHTTPHeaders(traceHeaders); err == nil {
			span := tracing.StartSpan("message_handler", opentracing.ChildOf(spanCtx))
			ctx = tracing.ContextWithSpan(ctx, span)
			defer span.Finish()
		}
	}

	span, ctx := tracing.TraceMessage(ctx, "consume", bus.config.RabbitMQ.Exchange, eventType)
	defer span.Finish()

	bus.logger.Info("Processando mensagem",
		zap.String("event_type", eventType),
		zap.String("event_id", eventID),
		zap.String("tenant_id", tenantID),
	)

	// Desserializar evento
	event, err := bus.serializer.Deserialize(delivery.Body, eventType)
	if err != nil {
		tracing.SetSpanError(span, err)
		bus.logger.Error("Erro ao desserializar evento", zap.Error(err))
		delivery.Nack(false, false) // Enviar para DLQ
		return
	}

	// Buscar handlers
	bus.mu.RLock()
	handlers, exists := bus.handlers[eventType]
	bus.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		bus.logger.Debug("Nenhum handler encontrado para evento",
			zap.String("event_type", eventType),
		)
		delivery.Ack(false)
		return
	}

	// Processar com cada handler
	success := true
	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			success = false
			tracing.SetSpanError(span, err)
			logging.LogError(ctx, bus.logger, "Erro ao processar evento", err,
				zap.String("event_type", eventType),
				zap.String("event_id", eventID),
			)
		}
	}

	// Registrar métricas
	if bus.metrics != nil {
		bus.metrics.RecordMessageReceived(
			bus.config.RabbitMQ.Queue,
			tenantID,
			success,
			time.Since(start),
		)
	}

	// ACK/NACK baseado no sucesso
	if success {
		delivery.Ack(false)
		bus.logger.Debug("Mensagem processada com sucesso",
			zap.String("event_type", eventType),
			zap.String("event_id", eventID),
		)
	} else {
		// Verificar se deve rejeitar ou tentar novamente
		retryCount := 0
		if count, ok := delivery.Headers["x-retry-count"].(int32); ok {
			retryCount = int(count)
		}

		maxRetries := 3
		if retryCount < maxRetries {
			// Republicar com incremento no contador de retry
			bus.retryMessage(ctx, delivery, retryCount+1)
			delivery.Ack(false)
		} else {
			// Enviar para DLQ após esgotar tentativas
			delivery.Nack(false, false)
		}
	}
}

// retryMessage republica mensagem com delay
func (bus *RabbitMQEventBus) retryMessage(ctx context.Context, delivery amqp.Delivery, retryCount int) {
	delay := time.Duration(retryCount*retryCount) * time.Second // Backoff exponencial

	go func() {
		time.Sleep(delay)

		headers := delivery.Headers
		if headers == nil {
			headers = make(amqp.Table)
		}
		headers["x-retry-count"] = int32(retryCount)

		err := bus.channel.Publish(
			bus.config.RabbitMQ.Exchange,
			delivery.RoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  delivery.ContentType,
				DeliveryMode: delivery.DeliveryMode,
				MessageId:    delivery.MessageId,
				Timestamp:    delivery.Timestamp,
				Headers:      headers,
				Body:         delivery.Body,
			},
		)

		if err != nil {
			logging.LogError(ctx, bus.logger, "Erro ao republicar mensagem", err,
				zap.String("message_id", delivery.MessageId),
				zap.Int("retry_count", retryCount),
			)
		}
	}()
}

// Stop para o event bus
func (bus *RabbitMQEventBus) Stop(ctx context.Context) error {
	close(bus.done)

	if bus.channel != nil {
		bus.channel.Close()
	}

	if bus.connection != nil {
		bus.connection.Close()
	}

	bus.logger.Info("Event bus RabbitMQ parado")
	return nil
}

// NewEventBus factory function para criar event bus baseado na configuração
func NewEventBus(cfg *config.Config, logger *zap.Logger, metrics *metrics.Metrics) (EventBus, error) {
	if cfg.IsDevelopment() {
		// Em desenvolvimento, usar event bus em memória para simplicidade
		return NewInMemoryEventBus(logger), nil
	}

	// Em produção, usar RabbitMQ
	return NewRabbitMQEventBus(cfg, logger, metrics)
}