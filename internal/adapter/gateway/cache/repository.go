// Package cache provides a two-level cache implementation for token storage.
// It uses in-memory cache as primary (L1) and Redis as secondary (L2).
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/port/repository"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// tokenEntry represents a cached token with its expiration time.
type tokenEntry struct {
	token     string
	expiresAt time.Time
}

// isValid checks if the token is still valid (not expired).
func (e *tokenEntry) isValid() bool {
	return time.Now().Before(e.expiresAt)
}

// CachedTokenRepository implements a two-level cache strategy.
// L1: In-memory cache (fast, volatile)
// L2: Redis cache (persistent, shared across instances)
type CachedTokenRepository struct {
	memory map[string]*tokenEntry
	redis  repository.TokenRepository
	mu     sync.RWMutex
}

// NewCachedTokenRepository creates a new CachedTokenRepository.
// If redis is nil, only in-memory cache will be used.
func NewCachedTokenRepository(redis repository.TokenRepository) *CachedTokenRepository {
	return &CachedTokenRepository{
		memory: make(map[string]*tokenEntry),
		redis:  redis,
	}
}

// SaveToken saves a token to both L1 (memory) and L2 (Redis) caches.
// This implements a write-through cache strategy.
func (r *CachedTokenRepository) SaveToken(ctx context.Context, key string, token string, ttlSeconds int) error {
	// Calculate expiration time (subtract 60 seconds for safety margin)
	ttl := time.Duration(ttlSeconds-60) * time.Second
	if ttl <= 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}
	expiresAt := time.Now().Add(ttl)

	// Save to L1 (in-memory) - always succeeds
	r.mu.Lock()
	r.memory[key] = &tokenEntry{
		token:     token,
		expiresAt: expiresAt,
	}
	r.mu.Unlock()
	logger.Debug("Cache", "Token saved to L1 (memory) for "+key)

	// Save to L2 (Redis) - best effort
	if r.redis != nil {
		if err := r.redis.SaveToken(ctx, key, token, ttlSeconds); err != nil {
			logger.Warning("Cache", "Failed to save token to L2 (Redis): "+err.Error())
			// Don't return error - L1 cache is sufficient
		} else {
			logger.Debug("Cache", "Token saved to L2 (Redis) for "+key)
		}
	}

	return nil
}

// GetToken retrieves a token from cache, checking L1 first, then L2.
func (r *CachedTokenRepository) GetToken(ctx context.Context, key string) (string, error) {
	// Check L1 (in-memory) first
	r.mu.RLock()
	if entry, ok := r.memory[key]; ok && entry.isValid() {
		r.mu.RUnlock()
		logger.Debug("Cache", "Token retrieved from L1 (memory) for "+key)
		return entry.token, nil
	}
	r.mu.RUnlock()

	// Check L2 (Redis) if available
	if r.redis != nil {
		token, err := r.redis.GetToken(ctx, key)
		if err == nil && token != "" {
			logger.Debug("Cache", "Token retrieved from L2 (Redis) for "+key)
			// Promote to L1 cache (use default TTL of 1 hour for promoted tokens)
			r.promoteToL1(key, token, 3600)
			return token, nil
		}
	}

	return "", nil
}

// IsTokenValid checks if a valid token exists in either cache level.
func (r *CachedTokenRepository) IsTokenValid(ctx context.Context, key string) bool {
	// Check L1 (in-memory) first
	r.mu.RLock()
	if entry, ok := r.memory[key]; ok && entry.isValid() {
		r.mu.RUnlock()
		return true
	}
	r.mu.RUnlock()

	// Check L2 (Redis) if available
	if r.redis != nil {
		return r.redis.IsTokenValid(ctx, key)
	}

	return false
}

// promoteToL1 promotes a token from L2 to L1 cache.
func (r *CachedTokenRepository) promoteToL1(key string, token string, ttlSeconds int) {
	ttl := time.Duration(ttlSeconds-60) * time.Second
	if ttl <= 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}

	r.mu.Lock()
	r.memory[key] = &tokenEntry{
		token:     token,
		expiresAt: time.Now().Add(ttl),
	}
	r.mu.Unlock()
	logger.Debug("Cache", "Token promoted to L1 (memory) for "+key)
}

// InvalidateToken removes a token from both L1 and L2 caches.
// This is used when an API returns an authentication error (401/400),
// indicating the cached token is no longer valid.
func (r *CachedTokenRepository) InvalidateToken(ctx context.Context, key string) error {
	// Remove from L1 (in-memory)
	r.mu.Lock()
	delete(r.memory, key)
	r.mu.Unlock()
	logger.Debug("Cache", "Token invalidated from L1 (memory) for "+key)

	// Remove from L2 (Redis) if available
	if r.redis != nil {
		if err := r.redis.InvalidateToken(ctx, key); err != nil {
			logger.Warning("Cache", "Failed to invalidate token from L2 (Redis): "+err.Error())
			// Don't return error - L1 invalidation is sufficient
		} else {
			logger.Debug("Cache", "Token invalidated from L2 (Redis) for "+key)
		}
	}

	return nil
}
