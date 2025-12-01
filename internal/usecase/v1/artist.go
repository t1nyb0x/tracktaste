// Package v1 contains V1 business logic for TrackTaste.
package v1

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

type ArtistUseCase struct {
	spotifyAPI external.SpotifyAPI
}

func NewArtistUseCase(spotifyAPI external.SpotifyAPI) *ArtistUseCase {
	return &ArtistUseCase{spotifyAPI: spotifyAPI}
}

func (uc *ArtistUseCase) FetchByID(ctx context.Context, artistID string) (*domain.Artist, error) {
	if artistID == "" {
		return nil, domain.ErrArtistNotFound
	}
	return uc.spotifyAPI.GetArtistByID(ctx, artistID)
}
