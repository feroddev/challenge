package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"github/feroddev/challengeV3/internal/pkg/crypto"
)

type Producer struct {
	conn      *nats.Conn
	logger    *zap.Logger
	config    ProducerConfig
	validator *SchemaValidator
	encryptor *crypto.Encryptor
}

type ProducerConfig struct {
	URL               string
	RetryAttempts     int
	RetryDelaySeconds int
	DeadLetterTopic   string
	EncryptionKey     string
}

func NewProducer(config ProducerConfig, logger *zap.Logger) (*Producer, error) {
	conn, err := nats.Connect(config.URL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao NATS: %w", err)
	}

	validator := NewSchemaValidator(logger)
	validator.RegisterDefaultSchemas()

	var encryptor *crypto.Encryptor
	if config.EncryptionKey != "" {
		encryptor, err = crypto.NewEncryptor(config.EncryptionKey)
		if err != nil {
			logger.Warn("Falha ao inicializar encriptador, dados nao serao criptografados", zap.Error(err))
		}
	}

	return &Producer{
		conn:      conn,
		logger:    logger,
		config:    config,
		validator: validator,
		encryptor: encryptor,
	}, nil
}

func (p *Producer) Publish(topic string, data interface{}) error {
	if p.validator != nil {
		err := p.validator.ValidateMessage(topic, data)
		if err != nil {
			p.logger.Error("Validacao de schema falhou", 
				zap.String("topic", topic), 
				zap.Error(err))
			return fmt.Errorf("erro de validacao de schema: %w", err)
		}
	}

	dataToSend := data
	if p.encryptor != nil {
		encryptedData, err := p.encryptSensitiveData(data)
		if err != nil {
			p.logger.Warn("Falha ao criptografar dados sensíveis", zap.Error(err))
		} else {
			dataToSend = encryptedData
		}
	}

	payload, err := json.Marshal(dataToSend)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	p.logger.Info("Publicando mensagem",
		zap.String("topic", topic),
		zap.Int("payload_size", len(payload)),
		zap.Time("timestamp", time.Now()))

	err = p.conn.Publish(topic, payload)
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem: %w", err)
	}

	p.logger.Info("Mensagem publicada com sucesso",
		zap.String("topic", topic),
		zap.Time("timestamp", time.Now()))

	return nil
}

func (p *Producer) PublishWithRetry(ctx context.Context, topic string, data interface{}, retries int, delay time.Duration) error {
	var lastErr error
	
	attempts := retries
	if attempts <= 0 {
		attempts = p.config.RetryAttempts
	}
	
	retryDelay := delay
	if retryDelay <= 0 {
		retryDelay = time.Duration(p.config.RetryDelaySeconds) * time.Second
	}
	
	for attempt := 0; attempt < attempts; attempt++ {
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
		case <-time.After(retryDelay):
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

func (p *Producer) Close() error {
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

type DeadLetterMessage struct {
	OriginalTopic string      `json:"original_topic"`
	Data          interface{} `json:"data"`
	Error         string      `json:"error"`
	Timestamp     time.Time   `json:"timestamp"`
}
