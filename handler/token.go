package handler

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/httpclient"
)

func SpotifyTokenHandler(w http.ResponseWriter) string {
    accessToken, err := httpclient.GetSpotifyAccessToken()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting access token: %v", err), http.StatusInternalServerError)
        return ""
    }

    return accessToken
}

func KKBoxTokenHandler(w http.ResponseWriter) string {
    accessToken, err := httpclient.GetKKBoxAccessToken()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting access token: %v", err), http.StatusInternalServerError)
        return ""
    }

    return accessToken
}