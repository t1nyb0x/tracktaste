package domain

type Track struct {
	Album             Album             `json:"album"`
	Artists           []Artist          `json:"artists"`
	AvailableMarkets  []string          `json:"available_markets"`
	DiscNumber        int16             `json:"disc_number"`
	DurationMs        int32             `json:"duration_ms"`
	Explicit          bool              `json:"explicit"`
	ExternalIDs       map[string]string `json:"external_ids"`
	ExternalURLs      map[string]string `json:"external_urls"`
	Href              string            `json:"href"`
	ID                string            `json:"id"`
	IsPlayable        bool              `json:"is_playable"`
	LinkedFrom        map[string]string `json:"linked_from"`
	Name              string            `json:"name"`
	Popularity        int16             `json:"popularity"`
	PreviewURL        string            `json:"preview_url"`
	TrackNumber       int16             `json:"track_number"`
	Typee             string            `json:"type"`
	URI               string            `json:"uri"`
	IsLocal           bool              `json:"is_local"`
}

type Artist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Typee        string            `json:"type"`
	URI          string            `json:"uri"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Album struct {
	AlbumType            string            `json:"album_type"`
	TotalTracks          int16             `json:"total_tracks"`
	AvailableMarkets     []string          `json:"available_markets"`
	ExternalURLs         map[string]string `json:"external_urls"`
	Href                 string            `json:"href"`
	ID                   string            `json:"id"`
	Images               []Image           `json:"images"`
	Name                 string            `json:"name"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	Restrictions         map[string]string `json:"restrictions"`
	Typee                string            `json:"type"`
	URI                  string            `json:"uri"`
	Artists              []Artist          `json:"artists"`
}
