// Package api provides HTTP handlers for the TrackTaste REST API.
// It implements endpoints for track, artist, album, and similar tracks operations.
package api

import (
	"github.com/t1nyb0x/tracktaste/internal/infra/kkbox"
	"github.com/t1nyb0x/tracktaste/internal/infra/spotify"
	"github.com/t1nyb0x/tracktaste/internal/service"
)

// Handler holds the dependencies for API handlers.
// It provides access to services and external API clients.
type Handler struct {
	// Artist is the service for artist-related operations.
	Artist *service.ArtistService
	// Track is the service for track-related operations.
	Track *service.TrackService
	// SpotifyClient is the client for Spotify API calls.
	SpotifyClient *spotify.Client
	// KKBOXClient is the client for KKBOX API calls.
	KKBOXClient *kkbox.Client
}

// NewHandler creates a new Handler with the given dependencies.
func NewHandler(artist *service.ArtistService, track *service.TrackService, spotifyClient *spotify.Client, kkboxClient *kkbox.Client) *Handler {
	return &Handler{
		Artist:        artist,
		Track:         track,
		SpotifyClient: spotifyClient,
		KKBOXClient:   kkboxClient,
	}
}
