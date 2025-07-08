package handler

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/httpclient"
)

func FetchArtistHandler(w http.ResponseWriter, r *http.Request) {
    accessToken := TokenHandler(w)
    if accessToken == "" {
        return
    }

    body, err := httpclient.Fetch("https://api.spotify.com/v1/artists/1bY7QMGccPmba1f1frZ8Xb", accessToken)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(body))
}