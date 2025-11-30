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

	// Audio Features API
	GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error)
	GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error)

	// Recommendations API
	GetRecommendations(ctx context.Context, params RecommendationParams) ([]domain.Track, error)

	// Artist Genres API
	GetArtistGenres(ctx context.Context, artistID string) ([]string, error)
	GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error)
}

// RecommendationParams represents parameters for Spotify Recommendations API.
type RecommendationParams struct {
	SeedTracks  []string // Max 5 combined with SeedArtists and SeedGenres
	SeedArtists []string
	SeedGenres  []string
	Limit       int // Max 100, default 20

	// Target parameters (optional)
	TargetTempo        *float64
	TargetEnergy       *float64
	TargetValence      *float64
	TargetDanceability *float64
	TargetAcousticness *float64

	// Min/Max parameters (optional)
	MinTempo  *float64
	MaxTempo  *float64
	MinEnergy *float64
	MaxEnergy *float64
}
