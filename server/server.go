package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/t1nyb0x/tracktaste/internal/api"
)

type Options struct {
	Addr string
	Deps Deps
}

type Deps struct {
	Handler *api.Handler
}

func New(opts Options, deps Deps) *http.Server {

	// init router
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Timeout(15 * time.Second), middleware.Logger)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	

	// router
	r.Route("/v1", func(r chi.Router) {
		r.Post("/artists/fetch", deps.Handler.FetchArtist)
		r.Get("/tracks/search", deps.Handler.SearchTrack)
		// r.Get("/tracks/similar", h.GetSimilarTrack)
	})

	// server setting
	return &http.Server{
		Handler: r,
		Addr:    opts.Addr,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
}