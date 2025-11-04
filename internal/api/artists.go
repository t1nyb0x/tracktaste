package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/service"
)

type fetchArtistReq struct {
	Name string `json:"name"`
}
type fetchArtistRes struct {
	Item any `json:"item"` // domain.Artist をそのまま返す
}

func (h *Handler) FetchArtist(w http.ResponseWriter, r *http.Request) {
	var in fetchArtistReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	art, err := h.Artist.FetchByName(r.Context(), in.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArtistNotFound):
			http.Error(w, "artist not found", http.StatusNotFound)
		default:
			http.Error(w, "upstream error", http.StatusBadGateway)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fetchArtistRes{Item: art})
}
