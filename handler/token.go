package handler

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/httpclient"
)

func TokenHandler(w http.ResponseWriter) string {
    accessToken, err := httpclient.GetAccessToken()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting access token: %v", err), http.StatusInternalServerError)
        return ""
    }

    return accessToken
}