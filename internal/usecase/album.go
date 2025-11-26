package usecase

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

// AlbumUseCase handles album-related business logic.
type AlbumUseCase struct {
	spotifyAPI external.SpotifyAPI
}

// NewAlbumUseCase creates a new AlbumUseCase instance.
func NewAlbumUseCase(spotifyAPI external.SpotifyAPI) *AlbumUseCase {
	return &AlbumUseCase{
		spotifyAPI: spotifyAPI,
	}
}

// FetchByID fetches an album by its Spotify ID.
func (uc *AlbumUseCase) FetchByID(ctx context.Context, id string) (*domain.Album, error) {
	if id == "" {
		return nil, domain.ErrAlbumNotFound
	}
	return uc.spotifyAPI.GetAlbumByID(ctx, id)
}
