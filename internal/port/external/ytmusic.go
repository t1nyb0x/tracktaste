package external

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// YouTubeMusicAPI defines the interface for YouTube Music operations via sidecar.
type YouTubeMusicAPI interface {
	// GetSimilarTracks retrieves similar tracks for a given YouTube video ID.
	// Returns a list of recommended tracks based on YouTube Music's radio feature.
	GetSimilarTracks(ctx context.Context, videoID string, limit int) ([]domain.YTMusicTrack, error)

	// SearchTracks searches for tracks on YouTube Music.
	// Useful for finding video IDs when only artist/track name is available.
	SearchTracks(ctx context.Context, query string, limit int) ([]domain.YTMusicTrack, error)
}
