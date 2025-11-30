// Package redis provides Redis-based token storage implementation.
package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

var client *redis.Client

// Init initializes the Redis connection.
func Init() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// TokenRepository implements port/repository.TokenRepository using Redis.
type TokenRepository struct{}

// NewTokenRepository creates a new TokenRepository.
func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

// SaveToken saves a token to Redis with TTL.
func (r *TokenRepository) SaveToken(ctx context.Context, key string, token string, ttlSeconds int) error {
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	ttl := time.Duration(ttlSeconds-60) * time.Second
	if ttl <= 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}
	redisKey := fmt.Sprintf("token:%s", key)
	if err := client.Set(ctx, redisKey, token, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	logger.Debug("Redis", fmt.Sprintf("Token saved for %s with TTL %v", key, ttl))
	return nil
}

// GetToken retrieves a token from Redis.
func (r *TokenRepository) GetToken(ctx context.Context, key string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	redisKey := fmt.Sprintf("token:%s", key)
	token, err := client.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

// IsTokenValid checks if a valid token exists.
func (r *TokenRepository) IsTokenValid(ctx context.Context, key string) bool {
	if client == nil {
		return false
	}
	redisKey := fmt.Sprintf("token:%s", key)
	exists, err := client.Exists(ctx, redisKey).Result()
	return err == nil && exists > 0
}

// InvalidateToken removes a token from Redis.
// This is called when an API returns an authentication error (401/400).
func (r *TokenRepository) InvalidateToken(ctx context.Context, key string) error {
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	redisKey := fmt.Sprintf("token:%s", key)
	if err := client.Del(ctx, redisKey).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	logger.Debug("Redis", fmt.Sprintf("Token invalidated for %s", key))
	return nil
}
