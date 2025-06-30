package configs

import (
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"os"
	"strconv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	NATS        messaging.ProducerConfig
	Topics      TopicsConfig
	Tracing     TracingConfig
	Redis       RedisConfig
	AWS         AWSConfig
	Auth        AuthConfig
	RateLimit   RateLimitConfig
	Crypto      CryptoConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type TopicsConfig struct {
	Gyroscope string
	GPS       string
	Photo     string
}

type TracingConfig struct {
	Endpoint string
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

type AuthConfig struct {
	JWTSecret     string
	TokenDuration int
}

type RateLimitConfig struct {
	Requests int
	Period   int
}

type CryptoConfig struct {
	EncryptionKey string
}

func NewConfig() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "telemetry"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		NATS: messaging.ProducerConfig{
			URL:               getEnv("NATS_URL", "nats://localhost:4222"),
			RetryAttempts:     getEnvInt("NATS_RETRY_ATTEMPTS", 3),
			RetryDelaySeconds: getEnvInt("NATS_RETRY_DELAY", 1),
			DeadLetterTopic:   getEnv("NATS_DEAD_LETTER_TOPIC", "telemetry.dead_letter"),
			EncryptionKey:     getEnv("ENCRYPTION_KEY", ""),
		},
		Topics: TopicsConfig{
			Gyroscope: getEnv("TOPIC_GYROSCOPE", "telemetry.gyroscope"),
			GPS:       getEnv("TOPIC_GPS", "telemetry.gps"),
			Photo:     getEnv("TOPIC_PHOTO", "telemetry.photo"),
		},
		Tracing: TracingConfig{
			Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
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
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "telemetry_secret_key_change_in_production"),
			TokenDuration: getEnvInt("JWT_TOKEN_DURATION", 24),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Period:   getEnvInt("RATE_LIMIT_PERIOD", 60),
		},
		Crypto: CryptoConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
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

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
