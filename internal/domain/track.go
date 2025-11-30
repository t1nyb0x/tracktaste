// Package domain defines the core business entities for TrackTaste.
// These entities are independent of external services and represent
// the unified data model used across the application.
package domain

// Track represents a music track in the domain model.
// It is platform-agnostic and contains unified track information.
type Track struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ISRC        *string  `json:"isrc,omitempty"`
	URL         string   `json:"url"`
	Popularity  *int     `json:"popularity,omitempty"`
	DiscNumber  int      `json:"disc_number"`
	TrackNumber int      `json:"track_number"`
	DurationMs  int      `json:"duration_ms,omitempty"`
	Explicit    bool     `json:"explicit,omitempty"`
	Artists     []Artist `json:"artists"`
	Album       Album    `json:"album"`
}

// SimpleTrack represents a simplified track without full album details.
// Used in album track listings.
type SimpleTrack struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	TrackNumber int      `json:"track_number"`
	Artists     []Artist `json:"artists"`
}

// SimilarTrack represents a track returned from similar tracks search.
type SimilarTrack struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ISRC        *string `json:"isrc,omitempty"`
	UPC         *string `json:"upc,omitempty"`
	URL         string  `json:"url"`
	Popularity  *int    `json:"popularity,omitempty"`
	TrackNumber int     `json:"track_number"`
	DurationMs  int     `json:"duration_ms"`
	Explicit    bool    `json:"explicit"`
	Album       Album   `json:"album"`
}

// TrackSearchResult represents a list of tracks from search results.
type TrackSearchResult struct {
	Items []Track `json:"items"`
}

// SimilarTracksResult represents a list of similar tracks.
type SimilarTracksResult struct {
	Items []SimilarTrack `json:"items"`
}
