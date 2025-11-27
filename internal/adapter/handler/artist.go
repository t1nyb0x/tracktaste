package handler

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/usecase"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

type ArtistHandler struct {
	artistUC *usecase.ArtistUseCase
}

func NewArtistHandler(artistUC *usecase.ArtistUseCase) *ArtistHandler {
	return &ArtistHandler{artistUC: artistUC}
}

type artistResult struct {
	URL        string        `json:"url"`
	Followers  string        `json:"followers"`
	Genres     []string      `json:"genres"`
	ID         string        `json:"id"`
	Images     []imageResult `json:"images"`
	Name       string        `json:"name"`
	Popularity *int          `json:"popularity"`
}

func (h *ArtistHandler) FetchByURL(w http.ResponseWriter, r *http.Request) {
	logger.Info("ArtistFetch", "リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	artistID, err := extractSpotifyArtistID(rawURL)
	if err != nil {
		if e, ok := err.(*extractError); ok {
			logger.Warning("ArtistFetch", e.Message)
			badRequest(w, e.Message, e.Code)
			return
		}
		badRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	artist, err := h.artistUC.FetchByID(r.Context(), artistID)
	if err != nil {
		logger.Error("ArtistFetch", "Spotify API エラー: "+err.Error())
		serviceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	images := make([]imageResult, len(artist.Images))
	for i, img := range artist.Images {
		images[i] = imageResult{URL: img.URL, Height: img.Height, Width: img.Width}
	}

	followers := "0"
	if artist.Followers != nil {
		followers = fmt.Sprintf("%d", *artist.Followers)
	}

	result := artistResult{
		URL:        artist.URL,
		Followers:  followers,
		Genres:     artist.Genres,
		ID:         artist.ID,
		Images:     images,
		Name:       artist.Name,
		Popularity: artist.Popularity,
	}

	logger.Info("ArtistFetch", "リクエスト完了")
	success(w, result)
}
