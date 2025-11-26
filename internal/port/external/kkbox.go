package external

import (
	"context"
)

// KKBOXTrackInfo represents minimal track info from KKBOX.
type KKBOXTrackInfo struct {
	ID   string
	Name string
	ISRC string
}

// KKBOXAPI defines the interface for KKBOX API operations.
type KKBOXAPI interface {
	// SearchByISRC searches for tracks by ISRC code.
	SearchByISRC(ctx context.Context, isrc string) (*KKBOXTrackInfo, error)

	// GetRecommendedTracks gets recommended tracks for a given KKBOX track ID.
	GetRecommendedTracks(ctx context.Context, trackID string) ([]KKBOXTrackInfo, error)

	// GetTrackDetail gets detailed track information including ISRC.
	GetTrackDetail(ctx context.Context, trackID string) (*KKBOXTrackInfo, error)
}
