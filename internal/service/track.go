package service

import (
	"context"
	"errors"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/spotify"
	"github.com/t1nyb0x/tracktaste/internal/repository/spotify"
)

var ErrTrackNotFound = errors.New("track not found")

type TrackService struct {
	Repo spotify.TrackRepo
}

func NewTrackService(r spotify.TrackRepo) *TrackService {
	return &TrackService{Repo: r}
}

func (s *TrackService) FetchById(ctx context.Context, id string) (domain.Track, error) {
	if id == "" {
		return domain.Track{}, ErrTrackNotFound
	}
	track, err := s.Repo.GetInfo(ctx, id)
	if err != nil || track.ID == "" {
		return domain.Track{}, ErrTrackNotFound
	}
	return track, nil
}

func (s *TrackService) SearchByQuery(ctx context.Context, query string) ([]domain.Track, error) {
	if query == "" {
		return nil, ErrTrackNotFound
	}
	tracks, err := s.Repo.Search(ctx, query)
	if err != nil || len(tracks) == 0 {
		return nil, ErrTrackNotFound
	}
	return tracks, nil
}