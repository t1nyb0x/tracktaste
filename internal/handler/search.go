package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/t1nyb0x/tracktaste/internal/httpclient"
)

// qパラメータに入った名前に一致するトラックを検索する
func SearchTrack(w http.ResponseWriter, r *http.Request) {
    accessToken := SpotifyTokenHandler(w)
	if accessToken == "" {
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(query))
	raw, err := httpclient.GetTrackInfo(apiURL, accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// JSON整形
    var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}