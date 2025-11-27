package spotify

import "github.com/t1nyb0x/tracktaste/internal/domain"

type rawTrack struct {
	Album       rawAlbum          `json:"album"`
	Artists     []rawSimpleArtist `json:"artists"`
	DiscNumber  int               `json:"disc_number"`
	DurationMs  int               `json:"duration_ms"`
	Explicit    bool              `json:"explicit"`
	ExternalIDs struct {
		ISRC string `json:"isrc"`
	} `json:"external_ids"`
	ExternalURLs map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	TrackNumber  int               `json:"track_number"`
}

func (r *rawTrack) toDomain() *domain.Track {
	artists := make([]domain.Artist, len(r.Artists))
	for i, a := range r.Artists {
		artists[i] = domain.Artist{
			ID:   a.ID,
			Name: a.Name,
			URL:  a.ExternalURLs["spotify"],
		}
	}

	album := r.Album.toDomainSimple()

	track := &domain.Track{
		ID:          r.ID,
		Name:        r.Name,
		URL:         r.ExternalURLs["spotify"],
		DiscNumber:  r.DiscNumber,
		TrackNumber: r.TrackNumber,
		DurationMs:  r.DurationMs,
		Explicit:    r.Explicit,
		Artists:     artists,
		Album:       *album,
	}

	if r.Popularity > 0 {
		pop := r.Popularity
		track.Popularity = &pop
	}
	if r.ExternalIDs.ISRC != "" {
		isrc := r.ExternalIDs.ISRC
		track.ISRC = &isrc
	}

	return track
}

type rawSimpleArtist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
}

type rawImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type rawAlbum struct {
	ExternalURLs map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Images       []rawImage        `json:"images"`
	Name         string            `json:"name"`
	ReleaseDate  string            `json:"release_date"`
	Artists      []rawSimpleArtist `json:"artists"`
	TotalTracks  int               `json:"total_tracks"`
	Popularity   int               `json:"popularity"`
	ExternalIDs  struct {
		UPC string `json:"upc"`
	} `json:"external_ids"`
	Genres []string `json:"genres"`
	Tracks struct {
		Items []rawSimpleTrack `json:"items"`
	} `json:"tracks"`
}

func (r *rawAlbum) toDomainSimple() *domain.Album {
	images := make([]domain.Image, len(r.Images))
	for i, img := range r.Images {
		images[i] = domain.Image{URL: img.URL, Height: img.Height, Width: img.Width}
	}

	artists := make([]domain.Artist, len(r.Artists))
	for i, a := range r.Artists {
		artists[i] = domain.Artist{ID: a.ID, Name: a.Name, URL: a.ExternalURLs["spotify"]}
	}

	return &domain.Album{
		ID:          r.ID,
		Name:        r.Name,
		URL:         r.ExternalURLs["spotify"],
		ReleaseDate: r.ReleaseDate,
		Images:      images,
		Artists:     artists,
	}
}

func (r *rawAlbum) toDomain() *domain.Album {
	album := r.toDomainSimple()
	album.TotalTracks = r.TotalTracks
	album.Genres = r.Genres

	if r.Popularity > 0 {
		pop := r.Popularity
		album.Popularity = &pop
	}
	if r.ExternalIDs.UPC != "" {
		upc := r.ExternalIDs.UPC
		album.UPC = &upc
	}

	tracks := make([]domain.SimpleTrack, len(r.Tracks.Items))
	for i, t := range r.Tracks.Items {
		trackArtists := make([]domain.Artist, len(t.Artists))
		for j, a := range t.Artists {
			trackArtists[j] = domain.Artist{ID: a.ID, Name: a.Name, URL: a.ExternalURLs["spotify"]}
		}
		tracks[i] = domain.SimpleTrack{
			ID:          t.ID,
			Name:        t.Name,
			URL:         t.ExternalURLs["spotify"],
			TrackNumber: t.TrackNumber,
			Artists:     trackArtists,
		}
	}
	album.Tracks = tracks

	return album
}

type rawSimpleTrack struct {
	Artists      []rawSimpleArtist `json:"artists"`
	ExternalURLs map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	TrackNumber  int               `json:"track_number"`
}

type rawArtist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Followers    struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres     []string   `json:"genres"`
	ID         string     `json:"id"`
	Images     []rawImage `json:"images"`
	Name       string     `json:"name"`
	Popularity int        `json:"popularity"`
}

func (r *rawArtist) toDomain() *domain.Artist {
	images := make([]domain.Image, len(r.Images))
	for i, img := range r.Images {
		images[i] = domain.Image{URL: img.URL, Height: img.Height, Width: img.Width}
	}

	artist := &domain.Artist{
		ID:     r.ID,
		Name:   r.Name,
		URL:    r.ExternalURLs["spotify"],
		Genres: r.Genres,
		Images: images,
	}

	if r.Followers.Total > 0 {
		f := r.Followers.Total
		artist.Followers = &f
	}
	if r.Popularity > 0 {
		p := r.Popularity
		artist.Popularity = &p
	}

	return artist
}
