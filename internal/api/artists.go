package api

import (
	"fmt"
	"net/http"

	"github.com/t1nyb0x/tracktaste/internal/api/response"
	"github.com/t1nyb0x/tracktaste/internal/infra/spotify"
	"github.com/t1nyb0x/tracktaste/internal/util"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// ArtistResult represents the artist response format per spec.
// It contains artist metadata including followers, genres, and images.
type ArtistResult struct {
	URL        string        `json:"url"`
	Followers  string        `json:"followers"`
	Genres     []string      `json:"genres"`
	ID         string        `json:"id"`
	Images     []ImageResult `json:"images"`
	Name       string        `json:"name"`
	Popularity *int          `json:"popularity"`
}

// FetchArtist handles GET /artist/fetch requests.
// It extracts the artist ID from a Spotify URL and returns formatted artist information.
//
// Query Parameters:
//   - url: Spotify artist URL (required)
//
// Response:
//   - 200: Artist information in ArtistResult format
//   - 400: Invalid URL or missing parameter
//   - 503: Spotify API error
func (h *Handler) FetchArtist(w http.ResponseWriter, r *http.Request) {
	logger.Info("ArtistFetch", "Spotify API リクエスト開始")

	rawURL := r.URL.Query().Get("url")

	// Extract artist ID with error handling
	artistID, err := util.ExtractSpotifyArtistID(rawURL)
	if err != nil {
		if extractErr, ok := err.(*util.ExtractError); ok {
			logger.Warning("ArtistFetch", extractErr.Message)
			response.BadRequest(w, extractErr.Message, extractErr.Code)
			return
		}
		response.BadRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	artist, err := h.SpotifyClient.FetchArtistById(r.Context(), artistID)
	if err != nil {
		logger.Error("ArtistFetch", "Spotify API エラー: "+err.Error())
		response.ServiceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	// Convert to response format per spec
	result := ArtistResult{
		URL:       artist.ExternalURLs["spotify"],
		Followers: formatFollowers(artist.Followers.Total),
		Genres:    artist.Genres,
		ID:        artist.ID,
		Images:    convertArtistImages(artist.Images),
		Name:      artist.Name,
	}

	// Handle nullable popularity
	if artist.Popularity > 0 {
		result.Popularity = &artist.Popularity
	}

	logger.Info("ArtistFetch", "リクエスト完了")
	response.Success(w, result)
}

// formatFollowers formats follower count as a string.
func formatFollowers(count int) string {
	return fmt.Sprintf("%d", count)
}

// convertArtistImages converts spotify Image slice to ImageResult slice.
func convertArtistImages(images []spotify.Image) []ImageResult {
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
