package handler

import (
	"encoding/json"
	"net/http"

	"github.com/t1nyb0x/tracktaste/handler/musicstax"
)

func GetTrackInfo(w http.ResponseWriter, r *http.Request) {
	cf, ua, err := GetCF()

	if err != nil {
		http.Error(w, "Failed to get cf_clearance cookie", http.StatusInternalServerError)
		return
	}

	raw, err := musicstax.FetchJSON("29vY6gIKRje259YNZ7FyDb", cf, ua)
	if err != nil {
		http.Error(w, "Failed to fetch track info", http.StatusInternalServerError)
		return
	}

	var pretty map[string]any
	json.Unmarshal(raw, &pretty)
	out, _ := json.MarshalIndent(pretty, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}