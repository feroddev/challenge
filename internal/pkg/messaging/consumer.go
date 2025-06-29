package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type MessageHandler func(data []byte) error

type Consumer struct {
	conn            *nats.Conn
	logger          *zap.Logger
	config          ConsumerConfig
	subscriptions   []*nats.Subscription
	shutdownContext context.Context
	shutdown        context.CancelFunc
}

type ConsumerConfig struct {
	URL               string
	RetryAttempts     int
	RetryDelaySeconds int
	DeadLetterTopic   string
}

func NewConsumer(config ConsumerConfig, logger *zap.Logger) (*Consumer, error) {
	conn, err := nats.Connect(config.URL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao NATS: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		conn:            conn,
		logger:          logger,
		config:          config,
		subscriptions:   make([]*nats.Subscription, 0),
		shutdownContext: ctx,
		shutdown:        cancel,
	}, nil
}

func (c *Consumer) Subscribe(topic string, handler MessageHandler) error {
	sub, err := c.conn.Subscribe(topic, func(msg *nats.Msg) {
		c.logger.Info("Mensagem recebida",
			zap.String("topic", topic),
			zap.Int("payload_size", len(msg.Data)))

		err := c.processMessageWithRetry(topic, msg.Data, handler)
		if err != nil {
			c.logger.Error("Falha ao processar mensagem após todas as tentativas",
				zap.String("topic", topic),
				zap.Error(err))

			deadLetterPayload := DeadLetterMessage{
				OriginalTopic: topic,
				Data:          msg.Data,
				Error:         err.Error(),
				Timestamp:     time.Now(),
			}

			deadLetterData, jsonErr := json.Marshal(deadLetterPayload)
			if jsonErr != nil {
				c.logger.Error("Erro ao serializar mensagem para dead letter",
					zap.Error(jsonErr))
				return
			}

			publishErr := c.conn.Publish(c.config.DeadLetterTopic, deadLetterData)
			if publishErr != nil {
				c.logger.Error("Erro ao publicar na fila de dead letter",
					zap.Error(publishErr))
			}
		}
	})

	if err != nil {
		return fmt.Errorf("erro ao assinar tópico %s: %w", topic, err)
	}

	c.subscriptions = append(c.subscriptions, sub)
	c.logger.Info("Assinatura criada com sucesso", zap.String("topic", topic))
	return nil
}

func (c *Consumer) processMessageWithRetry(topic string, data []byte, handler MessageHandler) error {
	var lastErr error

	for attempt := 0; attempt < c.config.RetryAttempts; attempt++ {
		select {
		case <-c.shutdownContext.Done():
			return fmt.Errorf("consumidor está sendo encerrado")
		default:
			err := handler(data)
			if err == nil {
				return nil
			}

			lastErr = err
			c.logger.Warn("Falha ao processar mensagem, tentando novamente",
				zap.String("topic", topic),
				zap.Int("attempt", attempt+1),
				zap.Int("max_attempts", c.config.RetryAttempts),
				zap.Error(err))

			time.Sleep(time.Duration(c.config.RetryDelaySeconds) * time.Second)
		}
	}

	return lastErr
}

func (c *Consumer) Close() error {
	c.shutdown()

	for _, sub := range c.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			c.logger.Error("Erro ao cancelar assinatura", zap.Error(err))
		}
	}

	c.conn.Close()
	return nil
}
