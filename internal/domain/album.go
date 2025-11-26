package domain

// Album represents a music album in the domain model.
// It is platform-agnostic and contains unified album information.
type Album struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	ReleaseDate string        `json:"release_date"`
	UPC         *string       `json:"upc,omitempty"`
	Popularity  *int          `json:"popularity,omitempty"`
	TotalTracks int           `json:"total_tracks,omitempty"`
	Genres      []string      `json:"genres,omitempty"`
	Images      []Image       `json:"images"`
	Artists     []Artist      `json:"artists"`
	Tracks      []SimpleTrack `json:"tracks,omitempty"`
}
