package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"github/feroddev/challengeV3/internal/services"
)

type MockRekognitionClient struct {
	mock.Mock
}

func (m *MockRekognitionClient) CompareFaces(sourceImage []byte, targetImage []byte) (float32, error) {
	args := m.Called(sourceImage, targetImage)
	return args.Get(0).(float32), args.Error(1)
}

type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) Get(key string) (interface{}, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockRedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func TestProcessPhotoRecognition(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockProducer, err := newMockProducer()
	if err != nil {
		t.Skip("Erro ao criar mock do produtor NATS, pulando teste")
	}
	
	mockRekognition := new(MockRekognitionClient)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope",
		GPSTopic:       "gps",
		PhotoTopic:     "photo",
	}
	
	service := services.NewTelemetryServiceWithRecognition(
		mockRepo,
		mockProducer.producer,
		mockRekognition,
		mockCache,
		config,
		logger,
	)
	
	now := time.Now()
	deviceID := "test-device-123"
	
	photo1 := core.Photo{
		ID:        1,
		Photo:     []byte("test-photo-1"),
		Timestamp: now.Add(-time.Hour),
		DeviceID:  deviceID,
	}
	
	photo2 := core.Photo{
		ID:        2,
		Photo:     []byte("test-photo-2"),
		Timestamp: now,
		DeviceID:  deviceID,
	}
	
	photos := []core.Photo{photo1}
	
	mockRepo.On("GetPhotosByDeviceID", deviceID).Return(photos, nil)
	mockRekognition.On("CompareFaces", photo1.Photo, photo2.Photo).Return(float32(0.95), nil)
	mockRepo.On("UpdatePhotoRecognition", photo2.ID, true, float32(0.95)).Return(nil)
	mockCache.On("Get", mock.Anything).Return(nil, core.ErrCacheMiss)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	
	err = service.ProcessPhotoRecognition(photo2)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRekognition.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestProcessPhotoRecognitionWithCache(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockProducer, err := newMockProducer()
	if err != nil {
		t.Skip("Erro ao criar mock do produtor NATS, pulando teste")
	}
	
	mockRekognition := new(MockRekognitionClient)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope",
		GPSTopic:       "gps",
		PhotoTopic:     "photo",
	}
	
	service := services.NewTelemetryServiceWithRecognition(
		mockRepo,
		mockProducer.producer,
		mockRekognition,
		mockCache,
		config,
		logger,
	)
	
	now := time.Now()
	deviceID := "test-device-123"
	
	photo := core.Photo{
		ID:        2,
		Photo:     []byte("test-photo-2"),
		Timestamp: now,
		DeviceID:  deviceID,
	}
	
	cacheResult := &core.RecognitionResult{
		Recognized: true,
		Similarity: 0.95,
	}
	
	mockCache.On("Get", mock.Anything).Return(cacheResult, nil)
	mockRepo.On("UpdatePhotoRecognition", photo.ID, true, float32(0.95)).Return(nil)
	
	err = service.ProcessPhotoRecognition(photo)
	
	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRekognition.AssertNotCalled(t, "CompareFaces")
}

func TestProcessPhotoRecognitionNoMatch(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockProducer, err := newMockProducer()
	if err != nil {
		t.Skip("Erro ao criar mock do produtor NATS, pulando teste")
	}
	
	mockRekognition := new(MockRekognitionClient)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope",
		GPSTopic:       "gps",
		PhotoTopic:     "photo",
	}
	
	service := services.NewTelemetryServiceWithRecognition(
		mockRepo,
		mockProducer.producer,
		mockRekognition,
		mockCache,
		config,
		logger,
	)
	
	now := time.Now()
	deviceID := "test-device-123"
	
	photo1 := core.Photo{
		ID:        1,
		Photo:     []byte("test-photo-1"),
		Timestamp: now.Add(-time.Hour),
		DeviceID:  deviceID,
	}
	
	photo2 := core.Photo{
		ID:        2,
		Photo:     []byte("test-photo-2"),
		Timestamp: now,
		DeviceID:  deviceID,
	}
	
	photos := []core.Photo{photo1}
	
	mockRepo.On("GetPhotosByDeviceID", deviceID).Return(photos, nil)
	mockRekognition.On("CompareFaces", photo1.Photo, photo2.Photo).Return(float32(0.30), nil)
	mockRepo.On("UpdatePhotoRecognition", photo2.ID, false, float32(0.30)).Return(nil)
	mockCache.On("Get", mock.Anything).Return(nil, core.ErrCacheMiss)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	
	err = service.ProcessPhotoRecognition(photo2)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRekognition.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
