package api

import (
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/api/response"
	"github.com/t1nyb0x/tracktaste/internal/infra/spotify"
	"github.com/t1nyb0x/tracktaste/internal/util"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// AlbumResult represents the album response format per spec.
// It contains album metadata including images, artists, tracks, and identification codes.
type AlbumResult struct {
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Images      []ImageResult       `json:"images"`
	Name        string              `json:"name"`
	ReleaseDate string              `json:"release_date"`
	Artists     []AlbumArtistResult `json:"artists"`
	Tracks      AlbumTracksResult   `json:"tracks"`
	Popularity  *int                `json:"popularity"`
	UPC         *string             `json:"upc"`
	Genres      []string            `json:"genres"`
}

// AlbumArtistResult represents artist information within an album response.
type AlbumArtistResult struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// AlbumTracksResult represents the tracks collection within an album response.
type AlbumTracksResult struct {
	Items []AlbumTrackItem `json:"items"`
}

// AlbumTrackItem represents a single track within an album.
type AlbumTrackItem struct {
	Artists     []AlbumArtistResult `json:"artists"`
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	TrackNumber int                 `json:"track_number"`
}

// FetchAlbum handles GET /album/fetch requests.
// It extracts the album ID from a Spotify URL and returns formatted album information.
//
// Query Parameters:
//   - url: Spotify album URL (required)
//
// Response:
//   - 200: Album information in AlbumResult format
//   - 400: Invalid URL or missing parameter
//   - 503: Spotify API error
func (h *Handler) FetchAlbum(w http.ResponseWriter, r *http.Request) {
	logger.Info("AlbumFetch", "Spotify API リクエスト開始")

	rawURL := r.URL.Query().Get("url")

	// Extract album ID with error handling
	albumID, err := util.ExtractSpotifyAlbumID(rawURL)
	if err != nil {
		if extractErr, ok := err.(*util.ExtractError); ok {
			logger.Warning("AlbumFetch", extractErr.Message)
			response.BadRequest(w, extractErr.Message, extractErr.Code)
			return
		}
		response.BadRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	album, err := h.SpotifyClient.FetchAlbumById(r.Context(), albumID)
	if err != nil {
		logger.Error("AlbumFetch", "Spotify API エラー: "+err.Error())
		response.ServiceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	// Convert to response format per spec
	result := AlbumResult{
		URL:         album.ExternalURLs["spotify"],
		ID:          album.ID,
		Images:      convertAlbumImages(album.Images),
		Name:        album.Name,
		ReleaseDate: album.ReleaseDate,
		Artists:     convertAlbumSimpleArtists(album.Artists),
		Tracks:      convertAlbumTracks(album.Tracks),
		Genres:      album.Genres,
	}

	// Handle nullable fields
	if album.Popularity > 0 {
		result.Popularity = &album.Popularity
	}
	if album.ExternalIDs.UPC != "" {
		result.UPC = &album.ExternalIDs.UPC
	}

	logger.Info("AlbumFetch", "リクエスト完了")
	response.Success(w, result)
}

// convertAlbumImages converts spotify Image slice to ImageResult slice.
func convertAlbumImages(images []spotify.Image) []ImageResult {
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

// convertAlbumSimpleArtists converts spotify SimpleArtist slice to AlbumArtistResult slice.
func convertAlbumSimpleArtists(artists []spotify.SimpleArtist) []AlbumArtistResult {
	result := make([]AlbumArtistResult, len(artists))
	for i, a := range artists {
		result[i] = AlbumArtistResult{
			URL:  a.ExternalURLs["spotify"],
			Name: a.Name,
		}
	}
	return result
}

// convertAlbumTracks converts spotify TracksPage to AlbumTracksResult.
func convertAlbumTracks(tracks spotify.TracksPage) AlbumTracksResult {
	items := make([]AlbumTrackItem, len(tracks.Items))
	for i, t := range tracks.Items {
		items[i] = AlbumTrackItem{
			Artists:     convertAlbumTrackArtists(t.Artists),
			URL:         t.ExternalURLs["spotify"],
			ID:          t.ID,
			Name:        t.Name,
			TrackNumber: t.TrackNumber,
		}
	}
	return AlbumTracksResult{Items: items}
}

// convertAlbumTrackArtists converts spotify SimpleArtist slice to AlbumArtistResult slice for track artists.
func convertAlbumTrackArtists(artists []spotify.SimpleArtist) []AlbumArtistResult {
	result := make([]AlbumArtistResult, len(artists))
	for i, a := range artists {
		result[i] = AlbumArtistResult{
			URL:  a.ExternalURLs["spotify"],
			Name: a.Name,
		}
	}
	return result
}
