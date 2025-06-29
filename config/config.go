package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	Redis  RedisConfig
	AWS    AWSConfig
	NATS   NATSConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type NATSConfig struct {
	URL               string
	GyroscopeTopic    string
	GPSTopic          string
	PhotoTopic        string
	RetryAttempts     int
	RetryDelaySeconds int
	DeadLetterTopic   string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENV", "development"),
		},
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "telemetry"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},
		NATS: NATSConfig{
			URL:               getEnv("NATS_URL", "nats://localhost:4222"),
			GyroscopeTopic:    getEnv("NATS_GYROSCOPE_TOPIC", "gyroscope.telemetry"),
			GPSTopic:          getEnv("NATS_GPS_TOPIC", "gps.telemetry"),
			PhotoTopic:        getEnv("NATS_PHOTO_TOPIC", "photo.telemetry"),
			RetryAttempts:     getEnvAsInt("NATS_RETRY_ATTEMPTS", 3),
			RetryDelaySeconds: getEnvAsInt("NATS_RETRY_DELAY", 5),
			DeadLetterTopic:   getEnv("NATS_DEAD_LETTER_TOPIC", "telemetry.deadletter"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
