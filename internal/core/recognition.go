package core

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type RecognitionResult struct {
	Recognized bool
	Similarity float32
}

type TelemetryRepository interface {
	SaveGyroscopeData(data Gyroscope) error
	SaveGPSData(data GPS) error
	SavePhotoData(data Photo) error
	GetGyroscopeData() ([]Gyroscope, error)
	GetGPSData() ([]GPS, error)
	GetPhotoData() ([]Photo, error)
	GetPhotosByDeviceID(deviceID string) ([]Photo, error)
	UpdatePhotoRecognition(id uint, recognized bool, similarity float32) error
}

type Producer interface {
	Publish(topic string, data interface{}) error
	PublishWithRetry(ctx context.Context, topic string, data interface{}, retries int, delay time.Duration) error
	Close() error
}
