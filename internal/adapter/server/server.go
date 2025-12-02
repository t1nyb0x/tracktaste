package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/t1nyb0x/tracktaste/internal/adapter/handler"
)

type Config struct {
	Addr string
}

type Handlers struct {
	Track     *handler.TrackHandler
	Artist    *handler.ArtistHandler
	Album     *handler.AlbumHandler
	Recommend *handler.RecommendHandler
	Health    *handler.HealthHandler
}

func New(cfg Config, h Handlers) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Timeout(15*time.Second), middleware.Logger)
	r.Get("/healthz", h.Health.Check)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/track/fetch", h.Track.FetchByURL)
		r.Get("/track/search", h.Track.Search)
		r.Get("/track/similar", h.Track.FetchSimilar)
		r.Get("/artist/fetch", h.Artist.FetchByURL)
		r.Get("/album/fetch", h.Album.FetchByURL)
	})

	r.Route("/v2", func(r chi.Router) {
		r.Get("/track/recommend", h.Recommend.FetchRecommendations)
	})

	return &http.Server{
		Handler:      r,
		Addr:         cfg.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 35 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
