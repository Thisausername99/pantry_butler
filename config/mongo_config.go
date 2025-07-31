package config

import (
	"os"
	"strconv"
	"time"
)

type MongoDBConfig struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

type Config struct {
	MongoDB MongoDBConfig
	Server  ServerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

func Load() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:            getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:       getEnv("MONGODB_DATABASE", "pantry_butler"),
			ConnectTimeout: getDurationEnv("MONGODB_CONNECT_TIMEOUT", 10*time.Second),
			MaxPoolSize:    getUint64Env("MONGODB_MAX_POOL_SIZE", 100),
			MinPoolSize:    getUint64Env("MONGODB_MIN_POOL_SIZE", 10),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "27107"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getUint64Env(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseUint(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
