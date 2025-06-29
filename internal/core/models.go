package core

import (
	"time"

	"gorm.io/gorm"
)

type Gyroscope struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	X         float64   `json:"x" validate:"required"`
	Y         float64   `json:"y" validate:"required"`
	Z         float64   `json:"z" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type GPS struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Latitude  float64   `json:"latitude" validate:"required"`
	Longitude float64   `json:"longitude" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Photo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Photo     string    `json:"photo" validate:"required" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
