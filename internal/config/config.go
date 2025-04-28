package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	Port          string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	cfg := &Config{
		// Database connection URL
		DatabaseURL: os.Getenv("DATABASE_URL"),
		// Port for the server
		Port: os.Getenv("PORT"),
		// Redis connection details
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}

	// Ensure that the Redis DB is correctly parsed as an integer
	// Default to DB 0 if not specified
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		cfg.RedisDB = 0
	} else {
		var err error
		cfg.RedisDB, err = strconv.Atoi(redisDB)
		if err != nil {
			log.Fatalf("Invalid Redis DB value in .env: %v", err)
		}
	}

	// Check if required environment variables are present
	if cfg.DatabaseURL == "" || cfg.Port == "" || cfg.RedisAddr == "" {
		log.Fatal("Required environment variables missing")
	}

	return cfg
}
