// Package v1 contains V1 business logic for TrackTaste.
package v1

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

func (uc *TrackUseCase) FetchByID(ctx context.Context, trackID string) (*domain.Track, error) {
	if trackID == "" {
		return nil, domain.ErrTrackNotFound
	}
	return uc.spotifyAPI.GetTrackByID(ctx, trackID)
}

func (uc *TrackUseCase) Search(ctx context.Context, query string) ([]domain.Track, error) {
	if query == "" {
		return nil, domain.ErrEmptyQuery
	}
	return uc.spotifyAPI.SearchTracks(ctx, query)
}
