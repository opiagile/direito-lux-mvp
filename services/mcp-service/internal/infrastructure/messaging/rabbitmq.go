package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
)

// MessageHandler função para processar mensagens
type MessageHandler func(ctx context.Context, delivery amqp.Delivery) error

// RabbitMQConnection gerencia a conexão com RabbitMQ
type RabbitMQConnection struct {
	config     *config.Config
	logger     *zap.Logger
	metrics    *metrics.Metrics
	conn       *amqp.Connection
	channel    *amqp.Channel
	handlers   map[string]MessageHandler
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewRabbitMQConnection cria uma nova conexão com RabbitMQ
func NewRabbitMQConnection(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *metrics.Metrics,
) (*RabbitMQConnection, error) {
	// RabbitMQ sempre habilitado se URL fornecida
	if cfg.RabbitMQ.URL == "" {
		logger.Info("RabbitMQ URL não configurada, desabilitando")
		return &RabbitMQConnection{
			config:   cfg,
			logger:   logger,
			metrics:  metrics,
			handlers: make(map[string]MessageHandler),
		}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	rabbitmq := &RabbitMQConnection{
		config:   cfg,
		logger:   logger,
		metrics:  metrics,
		handlers: make(map[string]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}

	// Conectar ao RabbitMQ
	if err := rabbitmq.connect(); err != nil {
		cancel()
		return nil, fmt.Errorf("erro ao conectar com RabbitMQ: %w", err)
	}

	// Declarar exchanges e filas
	if err := rabbitmq.setupTopology(); err != nil {
		cancel()
		return nil, fmt.Errorf("erro ao configurar topologia: %w", err)
	}

	// Iniciar monitoramento de conexão
	go rabbitmq.monitorConnection()

	logger.Info("Conexão com RabbitMQ estabelecida",
		zap.String("url", cfg.RabbitMQ.URL),
	)

	return rabbitmq, nil
}

// connect estabelece conexão com RabbitMQ
func (r *RabbitMQConnection) connect() error {
	conn, err := amqp.Dial(r.config.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("erro ao conectar: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("erro ao criar canal: %w", err)
	}

	r.conn = conn
	r.channel = channel

	return nil
}

// setupTopology configura exchanges e filas
func (r *RabbitMQConnection) setupTopology() error {
	if r.channel == nil {
		return nil
	}

	// Declarar exchange para eventos MCP
	err := r.channel.ExchangeDeclare(
		"mcp.events",    // nome
		"topic",         // tipo
		true,            // durable
		false,           // auto-delete
		false,           // internal
		false,           // no-wait
		nil,             // argumentos
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar exchange mcp.events: %w", err)
	}

	// Declarar exchange para comandos MCP
	err = r.channel.ExchangeDeclare(
		"mcp.commands",  // nome
		"direct",        // tipo
		true,            // durable
		false,           // auto-delete
		false,           // internal
		false,           // no-wait
		nil,             // argumentos
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar exchange mcp.commands: %w", err)
	}

	// Declarar fila de dead letter
	_, err = r.channel.QueueDeclare(
		"mcp.dead-letter", // nome
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // argumentos
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar fila dead-letter: %w", err)
	}

	// Declarar filas para bot interactions
	botQueues := []string{
		"mcp.whatsapp.messages",
		"mcp.telegram.messages",
		"mcp.slack.messages",
		"mcp.tools.executions",
		"mcp.sessions.events",
	}

	for _, queueName := range botQueues {
		_, err = r.channel.QueueDeclare(
			queueName, // nome
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			amqp.Table{
				"x-dead-letter-exchange": "mcp.dead-letter",
				"x-message-ttl":          300000, // 5 minutos
			},
		)
		if err != nil {
			return fmt.Errorf("erro ao declarar fila %s: %w", queueName, err)
		}
	}

	r.logger.Info("Topologia RabbitMQ configurada com sucesso")
	return nil
}

// Publish publica uma mensagem
func (r *RabbitMQConnection) Publish(ctx context.Context, exchange, routingKey string, data []byte) error {
	if r.channel == nil {
		return fmt.Errorf("conexão não estabelecida")
	}

	start := time.Now()
	tenantID := logging.GetTenantID(ctx)
	traceID := logging.GetTraceID(ctx)

	// Criar headers
	headers := amqp.Table{
		"tenant_id": tenantID,
		"trace_id":  traceID,
		"timestamp": time.Now().UTC(),
	}

	// Publicar mensagem
	err := r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
			Body:         data,
		},
	)

	// Registrar métricas
	if r.metrics != nil {
		r.metrics.RecordMessageSent(exchange, routingKey, tenantID, err == nil)
	}

	if err != nil {
		r.logger.Error("Erro ao publicar mensagem",
			zap.String("exchange", exchange),
			zap.String("routing_key", routingKey),
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return fmt.Errorf("erro ao publicar: %w", err)
	}

	r.logger.Debug("Mensagem publicada",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
		zap.String("tenant_id", tenantID),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

// PublishMCPEvent publica evento MCP
func (r *RabbitMQConnection) PublishMCPEvent(ctx context.Context, eventType string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	return r.Publish(ctx, "mcp.events", eventType, payload)
}

// PublishBotMessage publica mensagem para bot
func (r *RabbitMQConnection) PublishBotMessage(ctx context.Context, botType string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	queueName := fmt.Sprintf("mcp.%s.messages", botType)
	return r.Publish(ctx, "", queueName, payload)
}

// PublishToolExecution publica execução de ferramenta
func (r *RabbitMQConnection) PublishToolExecution(ctx context.Context, execution interface{}) error {
	payload, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("erro ao serializar execução: %w", err)
	}

	return r.Publish(ctx, "", "mcp.tools.executions", payload)
}

// Subscribe subscreve a uma fila
func (r *RabbitMQConnection) Subscribe(queueName string, handler MessageHandler) error {
	if r.channel == nil {
		return fmt.Errorf("conexão não estabelecida")
	}

	r.handlers[queueName] = handler

	// Consumir mensagens
	msgs, err := r.channel.Consume(
		queueName, // fila
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("erro ao consumir fila %s: %w", queueName, err)
	}

	// Processar mensagens em goroutine
	go r.processMessages(queueName, msgs, handler)

	r.logger.Info("Subscrição configurada",
		zap.String("queue", queueName),
	)

	return nil
}

// processMessages processa mensagens de uma fila
func (r *RabbitMQConnection) processMessages(queueName string, msgs <-chan amqp.Delivery, handler MessageHandler) {
	for {
		select {
		case <-r.ctx.Done():
			return
		case delivery, ok := <-msgs:
			if !ok {
				r.logger.Warn("Canal de mensagens fechado", zap.String("queue", queueName))
				return
			}

			r.handleMessage(queueName, delivery, handler)
		}
	}
}

// handleMessage processa uma mensagem individual
func (r *RabbitMQConnection) handleMessage(queueName string, delivery amqp.Delivery, handler MessageHandler) {
	start := time.Now()

	// Extrair contexto dos headers
	ctx := context.Background()
	if tenantID, ok := delivery.Headers["tenant_id"].(string); ok {
		ctx = logging.WithTenantID(ctx, tenantID)
	}
	if traceID, ok := delivery.Headers["trace_id"].(string); ok {
		ctx = logging.WithTraceID(ctx, traceID)
	}

	contextLogger := logging.FromContext(ctx, r.logger)

	// Processar mensagem
	err := handler(ctx, delivery)
	processingTime := time.Since(start)

	// Registrar métricas
	if r.metrics != nil {
		tenantID := logging.GetTenantID(ctx)
		r.metrics.RecordMessageReceived(queueName, tenantID, err == nil, processingTime)
	}

	// ACK ou NACK
	if err != nil {
		contextLogger.Error("Erro ao processar mensagem",
			zap.String("queue", queueName),
			zap.Duration("processing_time", processingTime),
			zap.Error(err),
		)

		// NACK com requeue
		delivery.Nack(false, true)
	} else {
		contextLogger.Debug("Mensagem processada",
			zap.String("queue", queueName),
			zap.Duration("processing_time", processingTime),
		)

		// ACK
		delivery.Ack(false)
	}
}

// monitorConnection monitora a conexão e reconecta se necessário
func (r *RabbitMQConnection) monitorConnection() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-r.conn.NotifyClose(make(chan *amqp.Error)):
			r.logger.Warn("Conexão RabbitMQ perdida, tentando reconectar...")

			// Tentar reconectar
			for {
				select {
				case <-r.ctx.Done():
					return
				default:
					if err := r.connect(); err != nil {
						r.logger.Error("Erro ao reconectar", zap.Error(err))
						time.Sleep(5 * time.Second)
						continue
					}

					if err := r.setupTopology(); err != nil {
						r.logger.Error("Erro ao reconfigurar topologia", zap.Error(err))
						time.Sleep(5 * time.Second)
						continue
					}

					// Resubscrever filas
					for queueName, handler := range r.handlers {
						if err := r.Subscribe(queueName, handler); err != nil {
							r.logger.Error("Erro ao resubscrever fila",
								zap.String("queue", queueName),
								zap.Error(err),
							)
						}
					}

					r.logger.Info("Reconexão RabbitMQ realizada com sucesso")
					break
				}
			}
		}
	}
}

// Close fecha a conexão
func (r *RabbitMQConnection) Close() error {
	r.cancel()

	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}

// Health verifica a saúde da conexão
func (r *RabbitMQConnection) Health(ctx context.Context) error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("conexão RabbitMQ não estabelecida")
	}

	if r.channel == nil {
		return fmt.Errorf("canal RabbitMQ não estabelecido")
	}

	return nil
}