package usecase

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

type TrackUseCase struct {
	spotifyAPI external.SpotifyAPI
}

func NewTrackUseCase(spotifyAPI external.SpotifyAPI) *TrackUseCase {
	return &TrackUseCase{spotifyAPI: spotifyAPI}
}

func (uc *TrackUseCase) FetchByID(ctx context.Context, id string) (*domain.Track, error) {
	if id == "" {
		return nil, domain.ErrTrackNotFound
	}
	return uc.spotifyAPI.GetTrackByID(ctx, id)
}

func (uc *TrackUseCase) Search(ctx context.Context, query string) ([]domain.Track, error) {
	if query == "" {
		return nil, domain.ErrEmptyQuery
	}
	return uc.spotifyAPI.SearchTracks(ctx, query)
}
