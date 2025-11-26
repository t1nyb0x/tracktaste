package service

import (
	"context"
	"errors"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/lastfm"
	"github.com/t1nyb0x/tracktaste/internal/repository"
)

var	ErrArtistNotFound = errors.New("artist not found")

type ArtistService struct {
	Repo repository.ArtistRepo
}

func NewArtistService(r repository.ArtistRepo) *ArtistService {
	return &ArtistService{Repo: r}
}

// nameで取得。getInfoで無ければsearchの先頭を返す。
func (s *ArtistService) FetchByName(ctx context.Context, name string) (domain.Artist, error) {
	if name == "" {
		return domain.Artist{}, ErrArtistNotFound
	}
	art, err := s.Repo.GetInfo(ctx, name)
	if err == nil && art.Name != "" {
		return art, nil
	}
	// fallback: search
	art, err = s.Repo.SearchFirst(ctx, name)
	if err != nil || art.Name == "" {
		return domain.Artist{}, ErrArtistNotFound
	}
	return art, nil
}