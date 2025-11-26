package api

import (
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/api/response"
	domain "github.com/t1nyb0x/tracktaste/internal/domain/spotify"
	"github.com/t1nyb0x/tracktaste/internal/util"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// TrackResult represents the track response format per spec.
// It contains track metadata including album, artists, and identification codes.
type TrackResult struct {
	Album       TrackAlbumResult    `json:"album"`
	Artists     []TrackArtistResult `json:"artists"`
	DiscNumber  int                 `json:"disc_number"`
	Popularity  *int                `json:"popularity"`
	ISRC        *string             `json:"isrc"`
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	TrackNumber int                 `json:"track_number"`
}

// TrackAlbumResult represents album information within a track response.
type TrackAlbumResult struct {
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Images      []ImageResult       `json:"images"`
	Name        string              `json:"name"`
	ReleaseDate string              `json:"release_date"`
	Artists     []TrackArtistResult `json:"artists"`
}

// TrackArtistResult represents artist information within a track response.
type TrackArtistResult struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ImageResult represents an image with dimensions.
type ImageResult struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// FetchTrackByURL handles GET /track/fetch requests.
// It extracts the track ID from a Spotify URL and returns formatted track information.
//
// Query Parameters:
//   - url: Spotify track URL (required)
//
// Response:
//   - 200: Track information in TrackResult format
//   - 400: Invalid URL or missing parameter
//   - 503: Spotify API error
func (h *Handler) FetchTrackByURL(w http.ResponseWriter, r *http.Request) {
	logger.Info("TrackFetch", "Spotify API リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	
	// Extract track ID with error handling
	trackID, err := util.ExtractSpotifyTrackID(rawURL)
	if err != nil {
		if extractErr, ok := err.(*util.ExtractError); ok {
			logger.Warning("TrackFetch", extractErr.Message)
			response.BadRequest(w, extractErr.Message, extractErr.Code)
			return
		}
		response.BadRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	track, err := h.Track.FetchById(r.Context(), trackID)
	if err != nil {
		logger.Error("TrackFetch", "Spotify API エラー: "+err.Error())
		response.ServiceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	// Convert to response format per spec
	result := TrackResult{
		Album: TrackAlbumResult{
			URL:         track.Album.ExternalURLs["spotify"],
			ID:          track.Album.ID,
			Name:        track.Album.Name,
			ReleaseDate: track.Album.ReleaseDate,
			Images:      convertImages(track.Album.Images),
			Artists:     convertAlbumArtists(track.Album.Artists),
		},
		Artists:     convertTrackArtists(track.Artists),
		DiscNumber:  int(track.DiscNumber),
		URL:         track.TrackExternalUrls.SpotifyTrackURL,
		ID:          track.ID,
		Name:        track.Name,
		TrackNumber: int(track.TrackNumber),
	}

	// Handle nullable fields
	if track.Popularity > 0 {
		pop := int(track.Popularity)
		result.Popularity = &pop
	}
	if track.TrackExternalIDs.ISRC != "" {
		result.ISRC = &track.TrackExternalIDs.ISRC
	}

	logger.Info("TrackFetch", "リクエスト完了")
	response.Success(w, result)
}

// SearchTrack handles GET /track/search requests.
// It searches for tracks by query string.
//
// Query Parameters:
//   - q: Search query (required)
//
// Response:
//   - 200: List of matching tracks
//   - 400: Missing query parameter
//   - 503: Spotify API error
func (h *Handler) SearchTrack(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "検索クエリが入力されていません", "EMPTY_QUERY")
		return
	}

	tracks, err := h.Track.SearchByQuery(r.Context(), query)
	if err != nil {
		logger.Error("TrackSearch", "Spotify API エラー: "+err.Error())
		response.ServiceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	response.Success(w, tracks)
}

// convertImages converts domain Image slice to ImageResult slice.
func convertImages(images []domain.Image) []ImageResult {
	result := make([]ImageResult, len(images))
	for i, img := range images {
		result[i] = ImageResult{
			URL:    img.URL,
			Height: img.Height,
			Width:  img.Width,
		}
	}
	return result
}

// convertAlbumArtists converts domain Artist slice to TrackArtistResult slice for album artists.
func convertAlbumArtists(artists []domain.Artist) []TrackArtistResult {
	result := make([]TrackArtistResult, len(artists))
	for i, a := range artists {
		result[i] = TrackArtistResult{
			URL:  a.ArtistExternalURLs.SpotifyArtistURL,
			ID:   a.ID,
			Name: a.ArtistName,
		}
	}
	return result
}

// convertTrackArtists converts domain Artist slice to TrackArtistResult slice for track artists.
func convertTrackArtists(artists []domain.Artist) []TrackArtistResult {
	result := make([]TrackArtistResult, len(artists))
	for i, a := range artists {
		result[i] = TrackArtistResult{
			URL:  a.ArtistExternalURLs.SpotifyArtistURL,
			ID:   a.ID,
			Name: a.ArtistName,
		}
	}
	return result
}
