package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/EhsanSepehriNasab/polling-platform/internal/config"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis() {
	// Load configuration
	cfg := config.Load()

	// Redis connection settings
	redisAddr := cfg.RedisAddr
	redisPassword := cfg.RedisPassword

	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default Redis address
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis server address
		Password: redisPassword, // No password set
		DB:       0,             // Default DB
	})

	// Check if Redis is available
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

func GetCacheClient() *redis.Client {
	return rdb
}

// DeleteKeysByPattern deletes keys matching the given pattern.
func DeleteKeysByPattern(pattern string) error {
	var cursor uint64
	for {
		// Use SCAN to find keys matching the pattern
		keys, newCursor, err := rdb.Scan(context.Background(), cursor, pattern, 0).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		// Delete the found keys
		if len(keys) > 0 {
			_, err := rdb.Del(context.Background(), keys...).Result()
			if err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
			log.Printf("Deleted %d keys matching pattern: %s", len(keys), pattern)
		}

		// If the cursor is 0, the iteration is complete
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return nil
}
