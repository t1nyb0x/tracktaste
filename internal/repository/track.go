package repository

import (
	"context"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/spotify"
)

type TrackRepo interface {
	GetInfo(ctx context.Context, id string) (domain.Track, error)
	Search(ctx context.Context, query string) ([]domain.Track, error)
}
