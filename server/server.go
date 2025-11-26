// Package server provides the HTTP server configuration for TrackTaste.
// It sets up routing using go-chi and configures middleware.
package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/t1nyb0x/tracktaste/internal/api"
)

// Options contains server configuration options.
type Options struct {
	// Addr is the address to listen on (e.g., ":8080").
	Addr string
	// Deps contains the server dependencies.
	Deps Deps
}

// Deps contains the dependencies required by the server.
type Deps struct {
	// Handler is the API handler containing endpoint implementations.
	Handler *api.Handler
}

// New creates a new HTTP server with the configured routes and middleware.
//
// Routes:
//   - GET /healthz - Health check endpoint
//   - GET /v1/track/fetch - Fetch track by Spotify URL
//   - GET /v1/track/search - Search tracks by query
//   - GET /v1/track/similar - Get similar tracks
//   - GET /v1/artist/fetch - Fetch artist by Spotify URL
//   - GET /v1/album/fetch - Fetch album by Spotify URL
//
// Middleware:
//   - RequestID: Adds a unique request ID to each request
//   - Recoverer: Recovers from panics
//   - Timeout: 15 second timeout for requests
//   - Logger: Logs request details
func New(opts Options, deps Deps) *http.Server {

	// init router
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Timeout(15*time.Second), middleware.Logger)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

	// router
	r.Route("/v1", func(r chi.Router) {
		// Track endpoints
		r.Get("/track/fetch", deps.Handler.FetchTrackByURL)
		r.Get("/track/search", deps.Handler.SearchTrack)
		r.Get("/track/similar", deps.Handler.FetchSimilarTracks)

		// Artist endpoint
		r.Get("/artist/fetch", deps.Handler.FetchArtist)

		// Album endpoint
		r.Get("/album/fetch", deps.Handler.FetchAlbum)
	})

	// server setting
	return &http.Server{
		Handler:      r,
		Addr:         opts.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 35 * time.Second, // Increased for similar tracks endpoint
		IdleTimeout:  60 * time.Second,
	}
}