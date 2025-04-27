package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, reading environment variables directly")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" || cfg.Port == "" {
		log.Fatal("Required environment variables missing")
	}

	return cfg
}
