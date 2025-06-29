package services

import (
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/repositories"
)

type TelemetryService interface {
	SaveGyroscopeData(data core.Gyroscope) error
	SaveGPSData(data core.GPS) error
	SavePhotoData(data core.Photo) error
	GetGyroscopeData() ([]core.Gyroscope, error)
	GetGPSData() ([]core.GPS, error)
	GetPhotoData() ([]core.Photo, error)
}

type telemetryService struct {
	repository repositories.TelemetryRepository
}

func NewTelemetryService(repository repositories.TelemetryRepository) TelemetryService {
	return &telemetryService{
		repository: repository,
	}
}

func (s *telemetryService) SaveGyroscopeData(data core.Gyroscope) error {
	return s.repository.SaveGyroscopeData(data)
}

func (s *telemetryService) SaveGPSData(data core.GPS) error {
	return s.repository.SaveGPSData(data)
}

func (s *telemetryService) SavePhotoData(data core.Photo) error {
	return s.repository.SavePhotoData(data)
}

func (s *telemetryService) GetGyroscopeData() ([]core.Gyroscope, error) {
	return s.repository.GetGyroscopeData()
}

func (s *telemetryService) GetGPSData() ([]core.GPS, error) {
	return s.repository.GetGPSData()
}

func (s *telemetryService) GetPhotoData() ([]core.Photo, error) {
	return s.repository.GetPhotoData()
}
