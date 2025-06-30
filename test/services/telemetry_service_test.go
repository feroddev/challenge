package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
	"github/feroddev/challengeV3/test/mocks"
)



func TestSaveGyroscopeData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	gyroscopeData := core.Gyroscope{
		X:         10.5,
		Y:         -5.2,
		Z:         3.7,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGyroscopeData", gyroscopeData).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.GyroscopeTopic, gyroscopeData, mock.Anything, mock.Anything).Return(nil)

	err := service.SaveGyroscopeData(gyroscopeData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestSaveGPSData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	gpsData := core.GPS{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGPSData", gpsData).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.GPSTopic, gpsData, mock.Anything, mock.Anything).Return(nil)

	err := service.SaveGPSData(gpsData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestSavePhotoData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	photoData := core.Photo{
		Photo:     "base64_encoded_string",
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SavePhotoData", photoData).Return(nil)
	mockProducer.On("PublishWithRetry", mock.Anything, config.PhotoTopic, photoData, mock.Anything, mock.Anything).Return(nil)

	err := service.SavePhotoData(photoData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestGetGyroscopeData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	expectedData := []core.Gyroscope{
		{
			ID:        1,
			X:         10.5,
			Y:         -5.2,
			Z:         3.7,
			Timestamp: time.Now(),
			DeviceID:  "test-device",
		},
	}

	mockRepo.On("GetGyroscopeData").Return(expectedData, nil)

	data, err := service.GetGyroscopeData()

	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
	mockRepo.AssertExpectations(t)
}

func TestGetGPSData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	expectedData := []core.GPS{
		{
			ID:        1,
			Latitude:  -23.5505,
			Longitude: -46.6333,
			Timestamp: time.Now(),
			DeviceID:  "test-device",
		},
	}

	mockRepo.On("GetGPSData").Return(expectedData, nil)

	data, err := service.GetGPSData()

	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
	mockRepo.AssertExpectations(t)
}

func TestGetPhotoData(t *testing.T) {
	mockRepo := new(mocks.MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	mockProducer := new(mocks.MockProducer)
	
	service := services.NewTelemetryService(mockRepo, mockProducer, config, logger)

	expectedData := []core.Photo{
		{
			ID:        1,
			Photo:     "base64_encoded_string",
			Timestamp: time.Now(),
			DeviceID:  "test-device",
		},
	}

	mockRepo.On("GetPhotoData").Return(expectedData, nil)

	data, err := service.GetPhotoData()

	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
	mockRepo.AssertExpectations(t)
}
