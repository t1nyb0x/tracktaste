package domain

type Artist struct {
	Name        string `json:"name"`
	Mbid        string `json:"mbid,omitempty"`
	Url         string `json:"url,omitempty"`
	ImageSmall  string `json:"image_small,omitempty"`
	ImageLarge  string `json:"image_large,omitempty"`
	Listeners   int64  `json:"listeners,omitempty"`
	Playcount   int64  `json:"playcount,omitempty"`
}