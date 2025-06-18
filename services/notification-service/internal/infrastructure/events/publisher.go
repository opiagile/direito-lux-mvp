package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
)

// Publisher interface para publicação de eventos
type Publisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, events []DomainEvent) error
	Close() error
}

// DomainEvent representa um evento de domínio
type DomainEvent struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Source      string      `json:"source"`
	Subject     string      `json:"subject"`
	Time        time.Time   `json:"time"`
	Data        interface{} `json:"data"`
	TenantID    string      `json:"tenant_id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	TraceID     string      `json:"trace_id,omitempty"`
	Version     string      `json:"version"`
	ContentType string      `json:"content_type"`
}

// RabbitMQPublisher implementação do Publisher para RabbitMQ
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
	config  *config.RabbitMQConfig
}

// NewRabbitMQPublisher cria novo publisher RabbitMQ
func NewRabbitMQPublisher(cfg *config.RabbitMQConfig, logger *zap.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("erro ao criar canal RabbitMQ: %w", err)
	}

	// Declarar exchange de eventos
	err = channel.ExchangeDeclare(
		"notification.events", // nome
		"topic",              // tipo
		true,                 // durable
		false,                // auto-delete
		false,                // internal
		false,                // no-wait
		nil,                  // argumentos
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("erro ao declarar exchange: %w", err)
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		logger:  logger,
		config:  cfg,
	}, nil
}

// Publish publica um evento
func (p *RabbitMQPublisher) Publish(ctx context.Context, event DomainEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	routingKey := fmt.Sprintf("notification.%s", event.Type)

	return p.channel.Publish(
		"notification.events", // exchange
		routingKey,           // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    event.ID,
			Body:         body,
			Headers: amqp.Table{
				"tenant_id": event.TenantID,
				"trace_id":  event.TraceID,
				"version":   event.Version,
			},
		},
	)
}

// PublishBatch publica múltiplos eventos
func (p *RabbitMQPublisher) PublishBatch(ctx context.Context, events []DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return fmt.Errorf("erro ao publicar evento %s: %w", event.ID, err)
		}
	}
	return nil
}

// Close fecha conexões
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}