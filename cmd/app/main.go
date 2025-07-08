package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/t1nyb0x/tracktaste/httpclient"
)

func main() {
    r := mux.NewRouter()

    // route
    r.HandleFunc("/fetch-artist", FetchArtistHandler).Methods("GET")

    log.Println("サーバー起動中: http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("サーバー起動に失敗しました: %v", err)
    }
    log.Println("サーバー終了")
}

func TokenHandler(w http.ResponseWriter) string {
    accessToken, err := httpclient.GetAccessToken()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting access token: %v", err), http.StatusInternalServerError)
        return ""
    }

    return accessToken
}

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