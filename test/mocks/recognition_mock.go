package mocks

import (
	"time"
	
	"github.com/stretchr/testify/mock"
)

type MockAWSRekognition struct {
	mock.Mock
}

func (m *MockAWSRekognition) CompareFaces(sourceImage, targetImage string) (bool, float32, error) {
	args := m.Called(sourceImage, targetImage)
	return args.Bool(0), args.Get(1).(float32), args.Error(2)
}

type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisCache) Set(key, value string, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}
