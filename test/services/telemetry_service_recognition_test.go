package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
	"github/feroddev/challengeV3/test/mocks"
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

func TestProcessPhotoRecognition(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockAWS := new(MockAWSRekognition)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	service.SetAWSRekognition(mockAWS)
	service.SetRedisCache(mockCache)

	deviceID := "test-device"
	photo1 := core.Photo{
		ID:        1,
		Photo:     "base64_photo_1",
		DeviceID:  deviceID,
		Timestamp: time.Now().Add(-1 * time.Hour),
	}
	photo2 := core.Photo{
		ID:        2,
		Photo:     "base64_photo_2",
		DeviceID:  deviceID,
		Timestamp: time.Now(),
	}
	photos := []core.Photo{photo1, photo2}

	mockRepo.On("GetPhotosByDeviceID", deviceID).Return(photos, nil)
	mockCache.On("Get", "recognition:1:2").Return("", nil)
	mockAWS.On("CompareFaces", photo1.Photo, photo2.Photo).Return(true, float32(90.5), nil)
	mockCache.On("Set", "recognition:1:2", "true:90.5", mock.Anything).Return(nil)
	mockRepo.On("UpdatePhotoRecognition", uint(2), true, float32(90.5)).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.PhotoTopic, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := service.ProcessPhotoRecognition(context.Background(), photo2)

	assert.NoError(t, err)
	assert.True(t, result.Recognized)
	assert.Equal(t, float32(90.5), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockAWS.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestProcessPhotoRecognitionWithCache(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockAWS := new(MockAWSRekognition)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	service.SetAWSRekognition(mockAWS)
	service.SetRedisCache(mockCache)

	deviceID := "test-device"
	photo1 := core.Photo{
		ID:        1,
		Photo:     "base64_photo_1",
		DeviceID:  deviceID,
		Timestamp: time.Now().Add(-1 * time.Hour),
	}
	photo2 := core.Photo{
		ID:        2,
		Photo:     "base64_photo_2",
		DeviceID:  deviceID,
		Timestamp: time.Now(),
	}
	photos := []core.Photo{photo1, photo2}

	mockRepo.On("GetPhotosByDeviceID", deviceID).Return(photos, nil)
	mockCache.On("Get", "recognition:1:2").Return("true:95.5", nil)
	mockRepo.On("UpdatePhotoRecognition", uint(2), true, float32(95.5)).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.PhotoTopic, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := service.ProcessPhotoRecognition(context.Background(), photo2)

	assert.NoError(t, err)
	assert.True(t, result.Recognized)
	assert.Equal(t, float32(95.5), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockAWS.AssertNotCalled(t, "CompareFaces")
}

func TestProcessPhotoRecognitionNoMatch(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	mockAWS := new(MockAWSRekognition)
	mockCache := new(MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	service.SetAWSRekognition(mockAWS)
	service.SetRedisCache(mockCache)

	deviceID := "test-device"
	photo2 := core.Photo{
		ID:        2,
		Photo:     "base64_photo_2",
		DeviceID:  deviceID,
		Timestamp: time.Now(),
	}

	mockRepo.On("GetPhotosByDeviceID", deviceID).Return([]core.Photo{}, nil)
	mockRepo.On("UpdatePhotoRecognition", uint(2), false, float32(0)).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.PhotoTopic, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := service.ProcessPhotoRecognition(context.Background(), photo2)

	assert.NoError(t, err)
	assert.False(t, result.Recognized)
	assert.Equal(t, float32(0), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Get")
	mockAWS.AssertNotCalled(t, "CompareFaces")
}
