package handler

import (
	"fmt"
	"regexp"
	"strings"
)

type extractError struct {
	Code    string
	Message string
}

func (e *extractError) Error() string {
	return e.Message
}

func extractSpotifyID(rawURL string, resourceType string) (string, error) {
	if rawURL == "" {
		return "", &extractError{Code: "EMPTY_PARAM", Message: "URLが入力されていません"}
	}

	if !strings.Contains(rawURL, "spotify.com") {
		return "", &extractError{Code: "NOT_SPOTIFY_URL", Message: "SpotifyのURLを入力してください"}
	}

	resourceTypes := []string{"track", "artist", "album"}
	for _, rt := range resourceTypes {
		if rt != resourceType && strings.Contains(rawURL, "/"+rt+"/") {
			return "", &extractError{
				Code:    "DIFFERENT_SPOTIFY_URL",
				Message: fmt.Sprintf("%sのURLを入力してください", getResourceTypeName(resourceType)),
			}
		}
	}

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

	return "", &extractError{Code: "INVALID_URL", Message: "無効なURL形式です"}
}

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

func extractSpotifyTrackID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "track")
}

func extractSpotifyArtistID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "artist")
}

func extractSpotifyAlbumID(rawURL string) (string, error) {
	return extractSpotifyID(rawURL, "album")
}
