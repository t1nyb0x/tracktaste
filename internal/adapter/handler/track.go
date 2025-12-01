package handler

import (
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	usecasev1 "github.com/t1nyb0x/tracktaste/internal/usecase/v1"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

type TrackHandler struct {
	trackUC   *usecasev1.TrackUseCase
	similarUC *usecasev1.SimilarTracksUseCase
}

func NewTrackHandler(trackUC *usecasev1.TrackUseCase, similarUC *usecasev1.SimilarTracksUseCase) *TrackHandler {
	return &TrackHandler{trackUC: trackUC, similarUC: similarUC}
}

type trackResult struct {
	Album       trackAlbumResult    `json:"album"`
	Artists     []trackArtistResult `json:"artists"`
	DiscNumber  int                 `json:"disc_number"`
	Popularity  *int                `json:"popularity"`
	ISRC        *string             `json:"isrc"`
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	TrackNumber int                 `json:"track_number"`
	DurationMs  int                 `json:"duration_ms"`
	Explicit    bool                `json:"explicit"`
}

type trackAlbumResult struct {
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Images      []imageResult       `json:"images"`
	Name        string              `json:"name"`
	ReleaseDate string              `json:"release_date"`
	Artists     []trackArtistResult `json:"artists"`
}

type trackArtistResult struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type imageResult struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func (h *TrackHandler) FetchByURL(w http.ResponseWriter, r *http.Request) {
	logger.Info("TrackFetch", "リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	trackID, err := extractSpotifyTrackID(rawURL)
	if err != nil {
		if e, ok := err.(*extractError); ok {
			logger.Warning("TrackFetch", e.Message)
			badRequest(w, e.Message, e.Code)
			return
		}
		badRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	track, err := h.trackUC.FetchByID(r.Context(), trackID)
	if err != nil {
		logger.Error("TrackFetch", "Spotify API エラー: "+err.Error())
		serviceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	result := convertTrackToResult(track)
	logger.Info("TrackFetch", "リクエスト完了")
	success(w, result)
}

func (h *TrackHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		badRequest(w, "検索クエリが入力されていません", "EMPTY_QUERY")
		return
	}

	tracks, err := h.trackUC.Search(r.Context(), query)
	if err != nil {
		logger.Error("TrackSearch", "Spotify API エラー: "+err.Error())
		serviceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	results := make([]trackResult, len(tracks))
	for i, t := range tracks {
		results[i] = convertTrackToResult(&t)
	}
	success(w, trackSearchResponse{Items: results})
}

type trackSearchResponse struct {
	Items []trackResult `json:"items"`
}

func (h *TrackHandler) FetchSimilar(w http.ResponseWriter, r *http.Request) {
	logger.Info("TrackSimilar", "リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	trackID, err := extractSpotifyTrackID(rawURL)
	if err != nil {
		if e, ok := err.(*extractError); ok {
			logger.Warning("TrackSimilar", e.Message)
			badRequest(w, e.Message, e.Code)
			return
		}
		badRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	result, err := h.similarUC.FetchSimilar(r.Context(), trackID)
	if err != nil {
		switch err {
		case domain.ErrISRCNotFound:
			badRequest(w, "ISRCが見つかりませんでした", "ISRC_NOT_FOUND")
		case domain.ErrTrackNotFound:
			notFound(w, "KKBOXで曲が見つかりませんでした", "KKBOX_TRACK_NOT_FOUND")
		default:
			logger.Error("TrackSimilar", "API エラー: "+err.Error())
			serviceUnavailable(w, "APIで問題が発生しているようです", "SOMETHING_API_ERROR")
		}
		return
	}

	resp := similarTracksResponse{Items: convertSimilarTracks(result.Items)}
	logger.Info("TrackSimilar", "リクエスト完了")
	success(w, resp)
}

type similarTracksResponse struct {
	Items []similarTrackResult `json:"items"`
}

type similarTrackResult struct {
	Album       similarAlbumResult `json:"album"`
	ISRC        *string            `json:"isrc"`
	UPC         *string            `json:"upc"`
	URL         string             `json:"url"`
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Popularity  *int               `json:"popularity"`
	TrackNumber int                `json:"track_number"`
	DurationMs  int                `json:"duration_ms"`
	Explicit    bool               `json:"explicit"`
}

type similarAlbumResult struct {
	URL         string                `json:"url"`
	ID          string                `json:"id"`
	Images      []imageResult         `json:"images"`
	Name        string                `json:"name"`
	ReleaseDate string                `json:"release_date"`
	Artists     []similarArtistResult `json:"artists"`
}

type similarArtistResult struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

func convertTrackToResult(t *domain.Track) trackResult {
	images := make([]imageResult, len(t.Album.Images))
	for i, img := range t.Album.Images {
		images[i] = imageResult{URL: img.URL, Height: img.Height, Width: img.Width}
	}

	albumArtists := make([]trackArtistResult, len(t.Album.Artists))
	for i, a := range t.Album.Artists {
		albumArtists[i] = trackArtistResult{URL: a.URL, ID: a.ID, Name: a.Name}
	}

	artists := make([]trackArtistResult, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = trackArtistResult{URL: a.URL, ID: a.ID, Name: a.Name}
	}

	return trackResult{
		Album: trackAlbumResult{
			URL:         t.Album.URL,
			ID:          t.Album.ID,
			Images:      images,
			Name:        t.Album.Name,
			ReleaseDate: t.Album.ReleaseDate,
			Artists:     albumArtists,
		},
		Artists:     artists,
		DiscNumber:  t.DiscNumber,
		Popularity:  t.Popularity,
		ISRC:        t.ISRC,
		URL:         t.URL,
		ID:          t.ID,
		Name:        t.Name,
		TrackNumber: t.TrackNumber,
		DurationMs:  t.DurationMs,
		Explicit:    t.Explicit,
	}
}

func convertSimilarTracks(tracks []domain.SimilarTrack) []similarTrackResult {
	results := make([]similarTrackResult, len(tracks))
	for i, t := range tracks {
		images := make([]imageResult, len(t.Album.Images))
		for j, img := range t.Album.Images {
			images[j] = imageResult{URL: img.URL, Height: img.Height, Width: img.Width}
		}

		artists := make([]similarArtistResult, len(t.Album.Artists))
		for j, a := range t.Album.Artists {
			artists[j] = similarArtistResult{URL: a.URL, Name: a.Name, ID: a.ID}
		}

		results[i] = similarTrackResult{
			Album: similarAlbumResult{
				URL:         t.Album.URL,
				ID:          t.Album.ID,
				Images:      images,
				Name:        t.Album.Name,
				ReleaseDate: t.Album.ReleaseDate,
				Artists:     artists,
			},
			ISRC:        t.ISRC,
			UPC:         t.UPC,
			URL:         t.URL,
			ID:          t.ID,
			Name:        t.Name,
			Popularity:  t.Popularity,
			TrackNumber: t.TrackNumber,
			DurationMs:  t.DurationMs,
			Explicit:    t.Explicit,
		}
	}
	return results
}
