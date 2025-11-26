// Package redis provides Redis client functionality for token caching and management.
// It handles connection to Redis server and provides methods for storing/retrieving
// OAuth tokens with automatic TTL expiration.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// TokenData represents the structure of token data stored in Redis.
// It contains the access token and its expiration timestamp.
type TokenData struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// Init initializes the Redis client connection.
// It reads configuration from environment variables:
//   - REDIS_HOST: Redis server host (default: localhost)
//   - REDIS_PORT: Redis server port (default: 6379)
//   - REDIS_PASSWORD: Redis password (default: empty)
//   - REDIS_DB: Redis database number (default: 0)
//
// Returns an error if the connection fails.
func Init() error {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return client
}

// SaveToken saves an OAuth token to Redis with automatic TTL expiration.
// The token is stored with key format "token:{service}" and will be automatically
// deleted after the specified expiration time.
//
// Parameters:
//   - ctx: Context for the operation
//   - service: Service name (e.g., "spotify", "kkbox")
//   - token: The access token to store
//   - expiresIn: Token expiration time in seconds
//
// Returns an error if the save operation fails.
func SaveToken(ctx context.Context, service string, token string, expiresIn int) error {
	key := fmt.Sprintf("token:%s", service)
	ttl := time.Duration(expiresIn) * time.Second

	data := TokenData{
		AccessToken: token,
		ExpiresAt:   time.Now().Add(ttl).Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	return client.Set(ctx, key, jsonData, ttl).Err()
}

// GetToken retrieves an OAuth token from Redis.
// It looks up the token using key format "token:{service}".
//
// Parameters:
//   - ctx: Context for the operation
//   - service: Service name (e.g., "spotify", "kkbox")
//
// Returns the access token string or an error if not found or retrieval fails.
func GetToken(ctx context.Context, service string) (string, error) {
	key := fmt.Sprintf("token:%s", service)

	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	var data TokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return data.AccessToken, nil
}

// IsTokenValid checks if a valid token exists in Redis for the specified service.
// A token is considered valid if it exists in Redis (TTL not expired).
//
// Parameters:
//   - ctx: Context for the operation
//   - service: Service name (e.g., "spotify", "kkbox")
//
// Returns true if a valid token exists, false otherwise.
func IsTokenValid(ctx context.Context, service string) bool {
	key := fmt.Sprintf("token:%s", service)
	return client.Exists(ctx, key).Val() == 1
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
