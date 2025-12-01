// Package external defines the external API interfaces (ports) for TrackTaste.
package external

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// DeezerAPI defines the interface for Deezer API operations.
// Deezer provides BPM, Duration, and Gain information for tracks.
type DeezerAPI interface {
	// GetTrackByISRC searches for a track by ISRC and returns its features.
	GetTrackByISRC(ctx context.Context, isrc string) (*domain.DeezerTrack, error)

	// SearchTrack searches for a track by title and artist (fallback when ISRC is not available).
	SearchTrack(ctx context.Context, title, artist string) (*domain.DeezerTrack, error)

	// GetTracksByISRCBatch retrieves multiple tracks by their ISRCs.
	// Returns a map of ISRC -> DeezerTrack. ISRCs not found will be omitted from the result.
	GetTracksByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.DeezerTrack, error)
}
