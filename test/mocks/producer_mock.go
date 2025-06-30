package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockProducer implementa a interface core.Producer para testes
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Publish(topic string, data interface{}) error {
	args := m.Called(topic, data)
	return args.Error(0)
}

func (m *MockProducer) PublishWithRetry(ctx context.Context, topic string, data interface{}, retries int, delay time.Duration) error {
	args := m.Called(ctx, topic, data, retries, delay)
	return args.Error(0)
}

func (m *MockProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}
