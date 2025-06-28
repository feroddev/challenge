package config

import (
	"os"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
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
