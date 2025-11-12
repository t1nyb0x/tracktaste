package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/t1nyb0x/tracktaste/internal/service"
)

type fetchTrackReq struct {
	ID string `json:"id"`
}

type fetchTrackRes struct {
	Item any `json:"item"` // domain.Trackを返す
}

func (h *Handler) FetchTrackByURL(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] FetchTrackByURL handler called")
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "url parameter is required", http.StatusBadRequest)
		return
	}

	// Extract track ID from Spotify URL
	trackID := extractTrackID(rawURL)
	if trackID == "" {
		http.Error(w, "invalid Spotify track URL", http.StatusBadRequest)
		return
	}

	track, err := h.Track.FetchById(r.Context(), trackID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTrackNotFound):
			http.Error(w, "track not found", http.StatusNotFound)
		default:
			http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fetchTrackRes{Item: track})
}

func (h *Handler) SearchTrack(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "q parameter is required", http.StatusBadRequest)
		return
	}

	tracks, err := h.Track.SearchByQuery(r.Context(), query)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTrackNotFound):
			http.Error(w, "no tracks found", http.StatusNotFound)
		default:
			http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fetchTrackRes{Item: tracks})
}

func extractTrackID(rawURL string) string {
	// Handle both open.spotify.com and api.spotify.com URLs
	patterns := []string{
		`open\.spotify\.com/track/([a-zA-Z0-9]+)`,
		`api\.spotify\.com/v1/tracks/([a-zA-Z0-9]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(rawURL)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
