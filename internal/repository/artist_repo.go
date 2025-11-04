package repository

import (
	"context"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/lastfm"
)

type ArtistRepo interface {
	GetInfo(ctx context.Context, name string) (domain.Artist, error)
	SearchFirst(ctx context.Context, query string) (domain.Artist, error)
}
