package core

import "time"

type Gyroscope struct {
	X         float64   `json:"x" validate:"required"`
	Y         float64   `json:"y" validate:"required"`
	Z         float64   `json:"z" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
}

type GPS struct {
	Latitude  float64   `json:"latitude" validate:"required"`
	Longitude float64   `json:"longitude" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
}

type Photo struct {
	Photo     string    `json:"photo" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
