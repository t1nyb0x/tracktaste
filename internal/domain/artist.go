package domain

// Artist represents a music artist in the domain model.
// It is platform-agnostic and contains unified artist information.
type Artist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Followers  *int     `json:"followers,omitempty"`
	Popularity *int     `json:"popularity,omitempty"`
	Genres     []string `json:"genres,omitempty"`
	Images     []Image  `json:"images,omitempty"`
}

// SimpleArtist represents a simplified artist reference.
// Used in track and album listings.
type SimpleArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
