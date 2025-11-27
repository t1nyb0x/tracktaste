// Package external defines the external API interfaces (ports) for TrackTaste.
package external

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// SpotifyAPI defines the interface for Spotify API operations.
type SpotifyAPI interface {
	GetTrackByID(ctx context.Context, id string) (*domain.Track, error)
	GetArtistByID(ctx context.Context, id string) (*domain.Artist, error)
	GetAlbumByID(ctx context.Context, id string) (*domain.Album, error)
	SearchTracks(ctx context.Context, query string) ([]domain.Track, error)
	SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error)
}
