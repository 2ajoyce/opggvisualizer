// internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
)

var config *Config

type Config struct {
	SummonerID   string
	DatabasePath string
	APIServer    apiConfig
}

type apiConfig struct {
	Port string
}

func loadConfig() (*Config, error) {
	summonerID := os.Getenv("SUMMONER_ID")
	if summonerID == "" {
		return nil, fmt.Errorf("SUMMONER_ID is not set")
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "data.db"
	}

	return &Config{
		SummonerID:   summonerID,
		DatabasePath: databasePath,
		APIServer: apiConfig{
			Port: getEnv("API_PORT", "8080"),
		},
	}, nil
}

func GetConfig() *Config {
	if config == nil {
		cfg, err := loadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		config = cfg
	}
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
