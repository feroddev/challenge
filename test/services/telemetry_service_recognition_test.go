package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
	"github/feroddev/challengeV3/test/mocks"
)

func TestProcessPhotoRecognition(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	mockAWS := new(mocks.MockAWSRekognition)
	mockCache := new(mocks.MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.recognition",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	telemetryService := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	telemetryService.SetAWSRekognition(mockAWS)
	telemetryService.SetRedisCache(mockCache)
	
	photo1 := core.Photo{
		ID:       1,
		DeviceID: "device1",
		Photo:    "base64photo1",
	}
	
	photo2 := core.Photo{
		ID:       2,
		DeviceID: "device1",
		Photo:    "base64photo2",
	}
	
	mockRepo.On("GetPhotosByDeviceID", "device1").Return([]core.Photo{photo1}, nil)
	mockCache.On("Get", "recognition:1:2").Return("", nil)
	mockAWS.On("CompareFaces", photo1.Photo, photo2.Photo).Return(true, float32(80.0), nil)
	mockCache.On("Set", "recognition:1:2", "true:80.0", mock.Anything).Return(nil)
	mockRepo.On("UpdatePhotoRecognition", uint(1), true, float32(80.0)).Return(nil)
	mockProducer.On("Publish", "photo.recognition", mock.Anything).Return(nil)

	result, err := telemetryService.ProcessPhotoRecognition(photo1, photo2)

	assert.NoError(t, err)
	assert.True(t, result.Recognized)
	assert.Equal(t, float32(80.0), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockAWS.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockRepo.AssertCalled(t, "UpdatePhotoRecognition", uint(1), true, float32(80.0))
	mockProducer.AssertCalled(t, "Publish", "photo.recognition", mock.Anything)
}

func TestProcessPhotoRecognitionWithCache(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	mockAWS := new(mocks.MockAWSRekognition)
	mockCache := new(mocks.MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.recognition",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	telemetryService := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	telemetryService.SetAWSRekognition(mockAWS)
	telemetryService.SetRedisCache(mockCache)
	
	photo1 := core.Photo{
		ID:       1,
		DeviceID: "device1",
		Photo:    "base64photo1",
	}
	
	photo2 := core.Photo{
		ID:       2,
		DeviceID: "device1",
		Photo:    "base64photo2",
	}
	
	mockRepo.On("GetPhotosByDeviceID", "device1").Return([]core.Photo{photo1}, nil)
	mockCache.On("Get", "recognition:1:2").Return("true:80.0", nil)
	mockRepo.On("UpdatePhotoRecognition", uint(1), true, float32(80.0)).Return(nil)
	mockProducer.On("Publish", "photo.recognition", mock.Anything).Return(nil)

	result, err := telemetryService.ProcessPhotoRecognition(photo1, photo2)

	assert.NoError(t, err)
	assert.True(t, result.Recognized)
	assert.Equal(t, float32(80.0), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockAWS.AssertNotCalled(t, "CompareFaces")
}

func TestProcessPhotoRecognitionNoMatch(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	mockAWS := new(mocks.MockAWSRekognition)
	mockCache := new(mocks.MockRedisCache)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		PhotoTopic: "photo.recognition",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	telemetryService := services.NewTelemetryService(mockRepo, mockProducer, config, logger)
	telemetryService.SetAWSRekognition(mockAWS)
	telemetryService.SetRedisCache(mockCache)
	
	photo2 := core.Photo{
		ID:       2,
		DeviceID: "device1",
		Photo:    "base64photo2",
	}
	
	mockRepo.On("GetPhotosByDeviceID", "device1").Return([]core.Photo{}, nil)
	mockRepo.On("UpdatePhotoRecognition", uint(2), false, float32(0)).Return(nil)
	mockProducer.On("Publish", "photo.recognition", mock.Anything).Return(nil)

	result, err := telemetryService.ProcessPhotoRecognition(photo2, photo2)

	assert.NoError(t, err)
	assert.False(t, result.Recognized)
	assert.Equal(t, float32(0), result.Similarity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Get")
	mockAWS.AssertNotCalled(t, "CompareFaces")
}
