package usecase

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

// ArtistUseCase handles artist-related business logic.
type ArtistUseCase struct {
	spotifyAPI external.SpotifyAPI
}

// NewArtistUseCase creates a new ArtistUseCase instance.
func NewArtistUseCase(spotifyAPI external.SpotifyAPI) *ArtistUseCase {
	return &ArtistUseCase{
		spotifyAPI: spotifyAPI,
	}
}

// FetchByID fetches an artist by its Spotify ID.
func (uc *ArtistUseCase) FetchByID(ctx context.Context, id string) (*domain.Artist, error) {
	if id == "" {
		return nil, domain.ErrArtistNotFound
	}
	return uc.spotifyAPI.GetArtistByID(ctx, id)
}
