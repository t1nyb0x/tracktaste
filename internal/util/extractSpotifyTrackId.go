// Package util provides utility functions for the TrackTaste application.
// It includes URL extraction utilities for Spotify resources.
package util

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractError represents an error during URL extraction.
// It contains a machine-readable Code and a human-readable Message in Japanese.
//
// Error codes:
//   - EMPTY_PARAM: URL parameter is empty
//   - NOT_SPOTIFY_URL: URL is not a Spotify URL
//   - DIFFERENT_SPOTIFY_URL: URL is a Spotify URL but for a different resource type
//   - INVALID_URL: URL format is invalid
type ExtractError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ExtractError) Error() string {
	return e.Message
}

// ExtractTrackID extracts track ID from Spotify URL.
// Deprecated: Use ExtractSpotifyTrackID for proper error handling.
func ExtractTrackID(rawURL string) string {
	id, _ := ExtractSpotifyTrackID(rawURL)
	return id
}

// ExtractSpotifyTrackID extracts track ID from a Spotify track URL.
// Returns the track ID and nil error on success, or empty string and ExtractError on failure.
//
// Supported URL formats:
//   - https://open.spotify.com/track/{id}
//   - https://open.spotify.com/intl-xx/track/{id}
//   - https://api.spotify.com/v1/tracks/{id}
func ExtractSpotifyTrackID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "track")
}

// ExtractSpotifyArtistID extracts artist ID from a Spotify artist URL.
// Returns the artist ID and nil error on success, or empty string and ExtractError on failure.
//
// Supported URL formats:
//   - https://open.spotify.com/artist/{id}
//   - https://open.spotify.com/intl-xx/artist/{id}
//   - https://api.spotify.com/v1/artists/{id}
func ExtractSpotifyArtistID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "artist")
}

// ExtractSpotifyAlbumID extracts album ID from a Spotify album URL.
// Returns the album ID and nil error on success, or empty string and ExtractError on failure.
//
// Supported URL formats:
//   - https://open.spotify.com/album/{id}
//   - https://open.spotify.com/intl-xx/album/{id}
//   - https://api.spotify.com/v1/albums/{id}
func ExtractSpotifyAlbumID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "album")
}

// extractSpotifyID is the internal implementation for extracting Spotify resource IDs.
// resourceType should be "track", "artist", or "album".
func extractSpotifyID(rawURL string, resourceType string) (string, error) {
	if rawURL == "" {
		return "", &ExtractError{Code: "EMPTY_PARAM", Message: "URLが入力されていません"}
	}

	// Check if it's a Spotify URL
	if !strings.Contains(rawURL, "spotify.com") {
		return "", &ExtractError{Code: "NOT_SPOTIFY_URL", Message: "SpotifyのURLを入力してください"}
	}

	// Check if URL contains the correct resource type
	resourceTypes := []string{"track", "artist", "album"}
	containsOtherType := false
	for _, rt := range resourceTypes {
		if rt != resourceType && strings.Contains(rawURL, "/"+rt+"/") {
			containsOtherType = true
			break
		}
	}

	if containsOtherType {
		return "", &ExtractError{
			Code:    "DIFFERENT_SPOTIFY_URL",
			Message: fmt.Sprintf("%sのURLを入力してください", getResourceTypeName(resourceType)),
		}
	}

	// Extract ID using regex patterns
	patterns := []string{
		fmt.Sprintf(`open\.spotify\.com/(?:intl-[a-z]{2}/)?%s/([A-Za-z0-9]+)`, resourceType),
		fmt.Sprintf(`api\.spotify\.com/v1/%ss/([a-zA-Z0-9]+)`, resourceType),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(rawURL)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", &ExtractError{Code: "INVALID_URL", Message: "無効なURL形式です"}
}

// getResourceTypeName returns the display name for a resource type.
func getResourceTypeName(resourceType string) string {
	switch resourceType {
	case "track":
		return "Track"
	case "artist":
		return "Artist"
	case "album":
		return "Album"
	default:
		return resourceType
	}
}