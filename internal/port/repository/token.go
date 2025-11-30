// Package repository defines the repository interfaces (ports) for TrackTaste.
package repository

import "context"

// TokenRepository defines the interface for token storage operations.
type TokenRepository interface {
	SaveToken(ctx context.Context, key string, token string, ttlSeconds int) error
	GetToken(ctx context.Context, key string) (string, error)
	IsTokenValid(ctx context.Context, key string) bool
	// InvalidateToken removes a token from the cache (used when API returns auth error)
	InvalidateToken(ctx context.Context, key string) error
}
