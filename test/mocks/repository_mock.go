package mocks

import (
	"github.com/stretchr/testify/mock"
	
	"github/feroddev/challengeV3/internal/core"
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
