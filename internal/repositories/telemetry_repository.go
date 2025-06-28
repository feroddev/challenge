package repositories

import (
	"github/feroddev/challengeV3/internal/core"
)

type TelemetryRepository interface {
	SaveGyroscopeData(data core.Gyroscope) error
	SaveGPSData(data core.GPS) error
	SavePhotoData(data core.Photo) error
}

type InMemoryTelemetryRepository struct {
	gyroscopeData []core.Gyroscope
	gpsData       []core.GPS
	photoData     []core.Photo
}

func NewInMemoryTelemetryRepository() TelemetryRepository {
	return &InMemoryTelemetryRepository{
		gyroscopeData: []core.Gyroscope{},
		gpsData:       []core.GPS{},
		photoData:     []core.Photo{},
	}
}

func (r *InMemoryTelemetryRepository) SaveGyroscopeData(data core.Gyroscope) error {
	r.gyroscopeData = append(r.gyroscopeData, data)
	return nil
}

func (r *InMemoryTelemetryRepository) SaveGPSData(data core.GPS) error {
	r.gpsData = append(r.gpsData, data)
	return nil
}

func (r *InMemoryTelemetryRepository) SavePhotoData(data core.Photo) error {
	r.photoData = append(r.photoData, data)
	return nil
}
