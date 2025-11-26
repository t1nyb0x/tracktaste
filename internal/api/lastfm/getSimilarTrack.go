package lastfm

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func GetSimilarTrack(artist string, track string, key string, limit string, w http.ResponseWriter) {
	base := "https://ws.audioscrobbler.com/2.0"
	v := url.Values{
		"method": {"track.getsimilar"},
		"artist": {artist},
		"track": {track},
		"api_key": {key},
		"format": {"json"},
		"limit": {limit}, // max 100
	}
	resp, err := http.Get(base + "?" + v.Encode())
	if err != nil {
		http.Error(w, "Failed to fetch similar tracks", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data struct {
		SimilarTracks struct {
			Track []struct {
				Name string `json:"name"`
				Artist struct {
					Name string
				} `json:"artist"`
			} `json:"track"`
		} `json:"similartracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}
	var titles []string
	for _, track := range data.SimilarTracks.Track {
		titles = append(titles, track.Name+" by "+track.Artist.Name)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(titles)
}