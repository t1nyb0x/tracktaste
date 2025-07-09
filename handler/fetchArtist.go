package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/httpclient"
)

type Artist struct {
    ID string `json:"artist_id"`
}

// FetchArtistHandler は、Spotify APIからアーティスト情報を取得するハンドラーです。
// このハンドラーは、事前にアクセストークンを取得
// し、指定されたアーティストの情報を取得してレスポンスとして返します。
// アクセストークンの取得に失敗した場合、エラーメッセージをHTTPレスポンスとして返します。
func FetchArtistHandler(w http.ResponseWriter, r *http.Request) {
    accessToken := TokenHandler(w)
    if accessToken == "" {
        return
    }

    var req Artist

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
        return
    }


    body, err := httpclient.GetArtistInfo("https://api.spotify.com/v1/artists/"+req.ID, accessToken)
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