package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Producer struct {
	conn   *nats.Conn
	logger *zap.Logger
	config ProducerConfig
}

type ProducerConfig struct {
	URL               string
	RetryAttempts     int
	RetryDelaySeconds int
	DeadLetterTopic   string
}

func NewProducer(config ProducerConfig, logger *zap.Logger) (*Producer, error) {
	conn, err := nats.Connect(config.URL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao NATS: %w", err)
	}

	return &Producer{
		conn:   conn,
		logger: logger,
		config: config,
	}, nil
}

func (p *Producer) Publish(topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	p.logger.Info("Publicando mensagem",
		zap.String("topic", topic),
		zap.Int("payload_size", len(payload)))

	err = p.conn.Publish(topic, payload)
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem: %w", err)
	}

	return nil
}

func (p *Producer) PublishWithRetry(ctx context.Context, topic string, data interface{}) error {
	var lastErr error
	
	for attempt := 0; attempt < p.config.RetryAttempts; attempt++ {
		err := p.Publish(topic, data)
		if err == nil {
			return nil
		}
		
		lastErr = err
		p.logger.Warn("Falha ao publicar mensagem, tentando novamente",
			zap.String("topic", topic),
			zap.Int("attempt", attempt+1),
			zap.Int("max_attempts", p.config.RetryAttempts),
			zap.Error(err))
		
		select {
		case <-time.After(time.Duration(p.config.RetryDelaySeconds) * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("contexto cancelado durante tentativas de publicação: %w", ctx.Err())
		}
	}
	
	p.logger.Error("Falha em todas as tentativas de publicação, enviando para dead letter",
		zap.String("topic", topic),
		zap.String("dead_letter_topic", p.config.DeadLetterTopic),
		zap.Error(lastErr))
	
	deadLetterPayload := DeadLetterMessage{
		OriginalTopic: topic,
		Data:          data,
		Error:         lastErr.Error(),
		Timestamp:     time.Now(),
	}
	
	err := p.Publish(p.config.DeadLetterTopic, deadLetterPayload)
	if err != nil {
		return fmt.Errorf("erro ao publicar na fila de dead letter: %w", err)
	}
	
	return lastErr
}

func (p *Producer) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

type DeadLetterMessage struct {
	OriginalTopic string      `json:"original_topic"`
	Data          interface{} `json:"data"`
	Error         string      `json:"error"`
	Timestamp     time.Time   `json:"timestamp"`
}
