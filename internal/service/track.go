package service

import (
	"context"
	"errors"
	"log"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/spotify"
	"github.com/t1nyb0x/tracktaste/internal/repository"
)

var ErrTrackNotFound = errors.New("track not found")

type TrackService struct {
	Repo repository.TrackRepo
}

func NewTrackService(r repository.TrackRepo) *TrackService {
	return &TrackService{Repo: r}
}

func (s *TrackService) FetchById(ctx context.Context, id string) (domain.Track, error) {
	log.Println("[DEBUG] FetchById called with id:", id, "at service layer")
	if id == "" {
		return domain.Track{}, ErrTrackNotFound
	}
	track, err := s.Repo.GetInfo(ctx, id)
	if err != nil {
		return domain.Track{}, ErrTrackNotFound
	}
	return track, nil
}

func (s *TrackService) SearchByQuery(ctx context.Context, query string) ([]domain.Track, error) {
	if query == "" {
		return nil, ErrTrackNotFound
	}
	tracks, err := s.Repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(tracks) == 0 {
		return nil, ErrTrackNotFound
	}
	return tracks, nil
}
