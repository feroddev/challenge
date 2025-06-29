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

func (m *MockTelemetryRepository) GetPhotosByDeviceID(deviceID string) ([]core.Photo, error) {
	args := m.Called(deviceID)
	return args.Get(0).([]core.Photo), args.Error(1)
}

func (m *MockTelemetryRepository) UpdatePhotoRecognition(id uint, recognized bool, similarity float32) error {
	args := m.Called(id, recognized, similarity)
	return args.Error(0)
}

type mockProducerWrapper struct {
	mock *mock.Mock
	producer *messaging.Producer
}

func newMockProducer() (*mockProducerWrapper, error) {
	logger, _ := zap.NewDevelopment()
	config := messaging.ProducerConfig{
		URL: "nats://localhost:4222",
	}
	
	producer, err := messaging.NewProducer(config, logger)
	if err != nil {
		return nil, err
	}
	
	m := &mock.Mock{}
	return &mockProducerWrapper{
		mock: m,
		producer: producer,
	}, nil
}

func TestSaveGyroscopeData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

	gyroscopeData := core.Gyroscope{
		X:         10.5,
		Y:         -5.2,
		Z:         3.7,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGyroscopeData", gyroscopeData).Return(nil)

	err = service.SaveGyroscopeData(gyroscopeData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSaveGPSData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

	gpsData := core.GPS{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SaveGPSData", gpsData).Return(nil)

	err = service.SaveGPSData(gpsData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSavePhotoData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{
		GyroscopeTopic: "gyroscope.telemetry",
		GPSTopic:       "gps.telemetry",
		PhotoTopic:     "photo.telemetry",
	}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

	photoData := core.Photo{
		Photo:     "base64_encoded_string",
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	mockRepo.On("SavePhotoData", photoData).Return(nil)

	err = service.SavePhotoData(photoData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetGyroscopeData(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

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
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

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
	mockRepo := new(MockTelemetryRepository)
	logger, _ := zap.NewDevelopment()
	
	config := services.TelemetryServiceConfig{}
	
	producer, err := messaging.NewProducer(messaging.ProducerConfig{
		URL: "nats://mock:4222",
	}, logger)
	if err != nil {
		t.Skip("Pulando teste: não foi possível criar o produtor NATS")
		return
	}
	
	service := services.NewTelemetryService(mockRepo, producer, config, logger)

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
