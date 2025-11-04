package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/service"
)

type fetchTrackReq struct {
	ID string `json:"id"`
}

type fetchTrackRes struct {
	Item any `json"item` // domain.Trackを返す
}

func (h *Handler) FetchTrack(w http.ResponseWriter, r *http.Request) {
	var in fetchTrackReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if in.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	track, err := h.Track.FetchById(r.Context(), in.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTrackNotFound):
			http.Error(w, "track not found", http.StatusNotFound)
		default:
			http.Error(w, "upstream error", http.StatusBadGateway)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fetchTrackRes{Item: track})
}