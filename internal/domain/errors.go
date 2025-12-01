package domain

import "errors"

// Domain errors represent business logic errors.
var (
	// ErrTrackNotFound indicates that a track was not found.
	ErrTrackNotFound = errors.New("track not found")

	// ErrArtistNotFound indicates that an artist was not found.
	ErrArtistNotFound = errors.New("artist not found")

	// ErrAlbumNotFound indicates that an album was not found.
	ErrAlbumNotFound = errors.New("album not found")

	// ErrISRCNotFound indicates that ISRC was not found for a track.
	ErrISRCNotFound = errors.New("ISRC not found")

	// ErrInvalidURL indicates that the provided URL is invalid.
	ErrInvalidURL = errors.New("invalid URL")

	// ErrEmptyQuery indicates that the search query is empty.
	ErrEmptyQuery = errors.New("empty query")

	// ErrExternalAPIError indicates an error from external API.
	ErrExternalAPIError = errors.New("external API error")

	// ErrTimeout indicates that the operation timed out.
	ErrTimeout = errors.New("operation timed out")

	// ErrNotFound is a generic not found error.
	ErrNotFound = errors.New("not found")
)

// ExtractError represents an error during URL extraction.
// It contains a machine-readable Code and a human-readable Message.
type ExtractError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ExtractError) Error() string {
	return e.Message
}

// Error codes for URL extraction errors.
const (
	ErrCodeEmptyParam          = "EMPTY_PARAM"
	ErrCodeNotSpotifyURL       = "NOT_SPOTIFY_URL"
	ErrCodeDifferentSpotifyURL = "DIFFERENT_SPOTIFY_URL"
	ErrCodeInvalidURL          = "INVALID_URL"
)
