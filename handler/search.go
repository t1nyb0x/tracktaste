package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/t1nyb0x/tracktaste/httpclient"
)


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

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(query))
	body, err := httpclient.GetArtistInfo(url, accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// JSON整形
    raw, _ :=json.Marshal(body)
    var pretty bytes.Buffer
    json.Indent(&pretty, raw, "", "  ")

    w.Header().Set("Content-Type", "application/json")
    w.Write(pretty.Bytes())
}