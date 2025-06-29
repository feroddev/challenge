package repositories

import (
	"github/feroddev/challengeV3/internal/core"
	"gorm.io/gorm"
)

type TelemetryRepository interface {
	SaveGyroscopeData(data core.Gyroscope) error
	SaveGPSData(data core.GPS) error
	SavePhotoData(data core.Photo) error
	GetGyroscopeData() ([]core.Gyroscope, error)
	GetGPSData() ([]core.GPS, error)
	GetPhotoData() ([]core.Photo, error)
	GetPhotosByDeviceID(deviceID string) ([]core.Photo, error)
	UpdatePhotoRecognition(id uint, recognized bool, similarity float32) error
}

type PostgresTelemetryRepository struct {
	db *gorm.DB
}

func NewPostgresTelemetryRepository(db *gorm.DB) TelemetryRepository {
	return &PostgresTelemetryRepository{
		db: db,
	}
}

func (r *PostgresTelemetryRepository) SaveGyroscopeData(data core.Gyroscope) error {
	return r.db.Create(&data).Error
}

func (r *PostgresTelemetryRepository) SaveGPSData(data core.GPS) error {
	return r.db.Create(&data).Error
}

func (r *PostgresTelemetryRepository) SavePhotoData(data core.Photo) error {
	return r.db.Create(&data).Error
}

func (r *PostgresTelemetryRepository) GetGyroscopeData() ([]core.Gyroscope, error) {
	var gyroscopeData []core.Gyroscope
	result := r.db.Find(&gyroscopeData)
	return gyroscopeData, result.Error
}

func (r *PostgresTelemetryRepository) GetGPSData() ([]core.GPS, error) {
	var gpsData []core.GPS
	result := r.db.Find(&gpsData)
	return gpsData, result.Error
}

func (r *PostgresTelemetryRepository) GetPhotoData() ([]core.Photo, error) {
	var photos []core.Photo
	result := r.db.Find(&photos)
	return photos, result.Error
}

func (r *PostgresTelemetryRepository) GetPhotosByDeviceID(deviceID string) ([]core.Photo, error) {
	var photos []core.Photo
	result := r.db.Where("device_id = ?", deviceID).Order("created_at DESC").Find(&photos)
	return photos, result.Error
}

func (r *PostgresTelemetryRepository) UpdatePhotoRecognition(id uint, recognized bool, similarity float32) error {
	result := r.db.Model(&core.Photo{}).Where("id = ?", id).Updates(map[string]interface{}{
		"recognized": recognized,
		"similarity": similarity,
	})
	return result.Error
}
