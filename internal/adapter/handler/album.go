package handler

import (
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/usecase"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

type AlbumHandler struct {
	albumUC *usecase.AlbumUseCase
}

func NewAlbumHandler(albumUC *usecase.AlbumUseCase) *AlbumHandler {
	return &AlbumHandler{albumUC: albumUC}
}

type albumResult struct {
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Images      []imageResult       `json:"images"`
	Name        string              `json:"name"`
	ReleaseDate string              `json:"release_date"`
	Artists     []albumArtistResult `json:"artists"`
	Tracks      albumTracksResult   `json:"tracks"`
	Popularity  *int                `json:"popularity"`
	UPC         *string             `json:"upc"`
	Genres      []string            `json:"genres"`
}

type albumArtistResult struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type albumTracksResult struct {
	Items []albumTrackItem `json:"items"`
}

type albumTrackItem struct {
	Artists     []albumArtistResult `json:"artists"`
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	TrackNumber int                 `json:"track_number"`
}

func (h *AlbumHandler) FetchByURL(w http.ResponseWriter, r *http.Request) {
	logger.Info("AlbumFetch", "リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	albumID, err := extractSpotifyAlbumID(rawURL)
	if err != nil {
		if e, ok := err.(*extractError); ok {
			logger.Warning("AlbumFetch", e.Message)
			badRequest(w, e.Message, e.Code)
			return
		}
		badRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	album, err := h.albumUC.FetchByID(r.Context(), albumID)
	if err != nil {
		logger.Error("AlbumFetch", "Spotify API エラー: "+err.Error())
		serviceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	images := make([]imageResult, len(album.Images))
	for i, img := range album.Images {
		images[i] = imageResult{URL: img.URL, Height: img.Height, Width: img.Width}
	}

	artists := make([]albumArtistResult, len(album.Artists))
	for i, a := range album.Artists {
		artists[i] = albumArtistResult{URL: a.URL, Name: a.Name}
	}

	trackItems := make([]albumTrackItem, len(album.Tracks))
	for i, t := range album.Tracks {
		trackArtists := make([]albumArtistResult, len(t.Artists))
		for j, a := range t.Artists {
			trackArtists[j] = albumArtistResult{URL: a.URL, Name: a.Name}
		}
		trackItems[i] = albumTrackItem{
			Artists:     trackArtists,
			URL:         t.URL,
			ID:          t.ID,
			Name:        t.Name,
			TrackNumber: t.TrackNumber,
		}
	}

	result := albumResult{
		URL:         album.URL,
		ID:          album.ID,
		Images:      images,
		Name:        album.Name,
		ReleaseDate: album.ReleaseDate,
		Artists:     artists,
		Tracks:      albumTracksResult{Items: trackItems},
		Popularity:  album.Popularity,
		UPC:         album.UPC,
		Genres:      album.Genres,
	}

	logger.Info("AlbumFetch", "リクエスト完了")
	success(w, result)
}
