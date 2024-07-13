package config

import (
	"fmt"
	"os"
	"time"
)

type HTTPServer struct {
	Address string
	Timeout time.Duration
}

type Config struct {
	Env      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	*HTTPServer
}

func LoadConfig() *Config {
	return &Config{
		Env:      getEnv("ENV", "local"),
		Host:     "db",
		Port:     getEnv("PG_PORT", "5432"),
		User:     getEnv("PG_USER", "user"),
		Password: getEnv("PG_PASSWORD", "user"),
		DBName:   getEnv("PG_DBNAME", "SCA"),
		HTTPServer: &HTTPServer{
			Address: getEnv("SERVER_ADDRESS", ":8787"),
			Timeout: parseTime("SERVER_TIMEOUT", "5s"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func parseTime(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		fmt.Printf("Error parsing %s duration: %s\n", key, err)
		duration, _ = time.ParseDuration(fallback)
	}
	return duration
}
