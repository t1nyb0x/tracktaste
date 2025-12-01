// Package v1 contains V1 business logic for TrackTaste.
package v1

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

type AlbumUseCase struct {
	spotifyAPI external.SpotifyAPI
}

func NewAlbumUseCase(spotifyAPI external.SpotifyAPI) *AlbumUseCase {
	return &AlbumUseCase{spotifyAPI: spotifyAPI}
}

func (uc *AlbumUseCase) FetchByID(ctx context.Context, albumID string) (*domain.Album, error) {
	if albumID == "" {
		return nil, domain.ErrAlbumNotFound
	}
	return uc.spotifyAPI.GetAlbumByID(ctx, albumID)
}
