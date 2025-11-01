package handler

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/httpclient"
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

func LastFMTokenHandler(w http.ResponseWriter) string {
    apiKey, err := httpclient.GetLastFMApiKey()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting API key: %v", err), http.StatusInternalServerError)
        return ""
    }

    return apiKey
}