// Package external defines the external API interfaces (ports) for TrackTaste.
package external

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// LastFMAPI defines the interface for Last.fm API operations.
// Last.fm provides similar tracks based on listening data.
type LastFMAPI interface {
	// GetSimilarTracks returns tracks similar to the given track.
	// Uses track.getSimilar API method.
	GetSimilarTracks(ctx context.Context, artist, track string, limit int) ([]domain.LastFMTrack, error)

	// GetSimilarTracksByMBID returns tracks similar to the given track using MusicBrainz ID.
	GetSimilarTracksByMBID(ctx context.Context, mbid string, limit int) ([]domain.LastFMTrack, error)
}
