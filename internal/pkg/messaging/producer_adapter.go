package messaging

import (
	"context"
	"time"

	"github/feroddev/challengeV3/internal/core"
)

// ProducerAdapter adapta o Producer NATS para implementar a interface core.Producer
type ProducerAdapter struct {
	producer *Producer
}

// NewProducerAdapter cria um novo adaptador para o Producer NATS
func NewProducerAdapter(producer *Producer) core.Producer {
	return &ProducerAdapter{
		producer: producer,
	}
}

// Publish publica uma mensagem no tópico especificado
func (a *ProducerAdapter) Publish(topic string, data interface{}) error {
	return a.producer.Publish(topic, data)
}

// PublishWithRetry publica uma mensagem com tentativas de reenvio
func (a *ProducerAdapter) PublishWithRetry(ctx context.Context, topic string, data interface{}, retries int, delay time.Duration) error {
	return a.producer.PublishWithRetry(ctx, topic, data, retries, delay)
}

// Close fecha a conexão com o NATS
func (a *ProducerAdapter) Close() error {
	return a.producer.Close()
}
