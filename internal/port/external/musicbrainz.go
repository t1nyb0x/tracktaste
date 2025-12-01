// Package external defines the external API interfaces (ports) for TrackTaste.
package external

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// MusicBrainzAPI defines the interface for MusicBrainz API operations.
// MusicBrainz provides tags and artist relation information.
type MusicBrainzAPI interface {
	// GetRecordingByISRC searches for a recording by ISRC.
	GetRecordingByISRC(ctx context.Context, isrc string) (*domain.MBRecording, error)

	// GetRecordingWithTags retrieves recording details including tags.
	GetRecordingWithTags(ctx context.Context, mbid string) (*domain.MBRecording, error)

	// GetArtistWithRelations retrieves artist details including tags and relations.
	GetArtistWithRelations(ctx context.Context, mbid string) (*domain.MBArtist, error)

	// GetRecordingsByISRCBatch retrieves multiple recordings by their ISRCs.
	// Returns a map of ISRC -> MBRecording. ISRCs not found will be omitted from the result.
	GetRecordingsByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.MBRecording, error)

	// GetArtistRecordings retrieves recordings by an artist (same artist's other tracks).
	// Returns recordings with ISRCs if available.
	GetArtistRecordings(ctx context.Context, artistMBID string, limit int) ([]domain.MBRecording, error)
}
