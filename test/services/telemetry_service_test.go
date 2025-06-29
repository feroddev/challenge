package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
)

type MockTelemetryRepository struct {
	mock.Mock
}

func (m *MockTelemetryRepository) SaveGyroscopeData(data core.Gyroscope) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockTelemetryRepository) SaveGPSData(data core.GPS) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockTelemetryRepository) SavePhotoData(data core.Photo) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockTelemetryRepository) GetGyroscopeData() ([]core.Gyroscope, error) {
	args := m.Called()
	return args.Get(0).([]core.Gyroscope), args.Error(1)
}

func (m *MockTelemetryRepository) GetGPSData() ([]core.GPS, error) {
	args := m.Called()
	return args.Get(0).([]core.GPS), args.Error(1)
}

func (m *MockTelemetryRepository) GetPhotoData() ([]core.Photo, error) {
	args := m.Called()
	return args.Get(0).([]core.Photo), args.Error(1)
}

func TestSaveGyroscopeData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := services.NewTelemetryService(mockRepo)

	gyroscopeData := core.Gyroscope{
		X:         10.5,
		Y:         -5.2,
		Z:         3.7,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGyroscopeData", gyroscopeData).Return(nil)

	err := service.SaveGyroscopeData(gyroscopeData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSaveGPSData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := services.NewTelemetryService(mockRepo)

	gpsData := core.GPS{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGPSData", gpsData).Return(nil)

	err := service.SaveGPSData(gpsData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSavePhotoData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := services.NewTelemetryService(mockRepo)

	photoData := core.Photo{
		Photo:     "base64_encoded_string",
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SavePhotoData", photoData).Return(nil)

	err := service.SavePhotoData(photoData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
