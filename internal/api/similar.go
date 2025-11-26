package api

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/api/response"
	"github.com/t1nyb0x/tracktaste/internal/infra/spotify"
	"github.com/t1nyb0x/tracktaste/internal/util"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// Concurrent processing configuration constants.
const (
	// maxConcurrent is the maximum number of concurrent Spotify API requests.
	maxConcurrent = 5
	// requestTimeout is the timeout for a single Spotify API request.
	requestTimeout = 5 * time.Second
	// overallTimeout is the overall timeout for the entire similar tracks request.
	overallTimeout = 30 * time.Second
	// maxSimilarTracks is the maximum number of similar tracks to return.
	maxSimilarTracks = 30
)

// SimilarTrackResult represents a single similar track in the response.
type SimilarTrackResult struct {
	Album       SimilarAlbumResult  `json:"album"`
	ISRC        *string             `json:"isrc"`
	UPC         *string             `json:"upc"`
	URL         string              `json:"url"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Popularity  *int                `json:"popularity"`
	TrackNumber int                 `json:"track_number"`
}

// SimilarAlbumResult represents album information within a similar track response.
type SimilarAlbumResult struct {
	URL         string                `json:"url"`
	ID          string                `json:"id"`
	Images      []ImageResult         `json:"images"`
	Name        string                `json:"name"`
	ReleaseDate string                `json:"release_date"`
	Artists     []SimilarArtistResult `json:"artists"`
}

// SimilarArtistResult represents artist information within a similar track response.
type SimilarArtistResult struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

// SimilarTracksResponse represents the response format for similar tracks endpoint.
type SimilarTracksResponse struct {
	Items []SimilarTrackResult `json:"items"`
}

// FetchSimilarTracks handles GET /track/similar requests.
// It finds similar tracks using KKBOX recommendations and returns them with Spotify metadata.
//
// Process Flow:
//  1. Extract track ID from Spotify URL
//  2. Get ISRC from Spotify track
//  3. Search KKBOX by ISRC
//  4. Get recommended tracks from KKBOX
//  5. Get track details with ISRC from KKBOX
//  6. Search Spotify by ISRC in parallel
//  7. Remove duplicates and sort by popularity
//
// Query Parameters:
//   - url: Spotify track URL (required)
//
// Response:
//   - 200: Similar tracks in SimilarTracksResponse format
//   - 400: Invalid URL, missing parameter, or ISRC not found
//   - 404: Track not found in KKBOX
//   - 503: Spotify or KKBOX API error
func (h *Handler) FetchSimilarTracks(w http.ResponseWriter, r *http.Request) {
	logger.Info("TrackSimilar", "リクエスト開始")

	// Create context with overall timeout
	ctx, cancel := context.WithTimeout(r.Context(), overallTimeout)
	defer cancel()

	rawURL := r.URL.Query().Get("url")

	// Extract track ID
	trackID, err := util.ExtractSpotifyTrackID(rawURL)
	if err != nil {
		if extractErr, ok := err.(*util.ExtractError); ok {
			logger.Warning("TrackSimilar", extractErr.Message)
			response.BadRequest(w, extractErr.Message, extractErr.Code)
			return
		}
		response.BadRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	// Step 1: Get track info from Spotify
	logger.Info("TrackSimilar", "Spotify APIからトラック情報を取得")
	var track spotify.SearchTrack
	if err := h.SpotifyClient.FetchById(ctx, trackID, &track); err != nil {
		logger.Error("TrackSimilar", "Spotify API エラー: "+err.Error())
		response.ServiceUnavailable(w, "Spotify APIで問題が発生しているようです", "SOMETHING_SPOTIFY_ERROR")
		return
	}

	// Step 2: Extract ISRC
	isrc := track.ExternalIDs.ISRC
	if isrc == "" {
		logger.Warning("TrackSimilar", "ISRCが見つかりません")
		response.BadRequest(w, "ISRCが見つかりませんでした", "ISRC_NOT_FOUND")
		return
	}

	// Step 3: Search KKBOX
	logger.Info("TrackSimilar", "KKBOXで検索")
	kkboxSearch, err := h.KKBOXClient.Search(ctx, isrc)
	if err != nil {
		logger.Error("TrackSimilar", "KKBOX API エラー: "+err.Error())
		response.ServiceUnavailable(w, "KKBOX APIで問題が発生しているようです", "SOMETHING_KKBOX_ERROR")
		return
	}

	if len(kkboxSearch.Tracks.Data) == 0 {
		logger.Warning("TrackSimilar", "KKBOXで曲が見つかりません")
		response.NotFound(w, "KKBOXで曲が見つかりませんでした", "KKBOX_TRACK_NOT_FOUND")
		return
	}

	// Step 4: Get KKBOX track ID
	kkboxTrackID := kkboxSearch.Tracks.Data[0].ID

	// Step 5: Get recommended tracks from KKBOX
	logger.Info("TrackSimilar", "KKBOXからレコメンドトラックを取得")
	recommended, err := h.KKBOXClient.GetRecommendedTracks(ctx, kkboxTrackID)
	if err != nil {
		logger.Error("TrackSimilar", "KKBOX API エラー: "+err.Error())
		response.ServiceUnavailable(w, "KKBOX APIで問題が発生しているようです", "SOMETHING_KKBOX_ERROR")
		return
	}

	if len(recommended.Data) == 0 {
		logger.Warning("TrackSimilar", "レコメンドトラックが見つかりません")
		response.Success(w, SimilarTracksResponse{Items: []SimilarTrackResult{}})
		return
	}

	// Step 6: Get track details from KKBOX to extract ISRC
	logger.Info("TrackSimilar", "KKBOXから詳細情報を取得")
	isrcList := make([]string, 0, len(recommended.Data))
	for _, track := range recommended.Data {
		detail, err := h.KKBOXClient.GetTrackDetail(ctx, track.ID)
		if err != nil {
			continue
		}
		if detail.ISRC != "" {
			isrcList = append(isrcList, detail.ISRC)
		}
	}

	// Step 7-8: Search Spotify with ISRC (parallel processing)
	logger.Info("TrackSimilar", "Spotifyで並列検索開始")
	similarTracks := h.searchSpotifyParallel(ctx, isrcList)

	// Step 9-10: Remove duplicates and limit to max tracks
	similarTracks = removeDuplicates(similarTracks)

	// Sort by popularity (descending)
	sort.Slice(similarTracks, func(i, j int) bool {
		popI := 0
		popJ := 0
		if similarTracks[i].Popularity != nil {
			popI = *similarTracks[i].Popularity
		}
		if similarTracks[j].Popularity != nil {
			popJ = *similarTracks[j].Popularity
		}
		return popI > popJ
	})

	// Limit to max tracks
	if len(similarTracks) > maxSimilarTracks {
		similarTracks = similarTracks[:maxSimilarTracks]
	}

	logger.Info("TrackSimilar", "リクエスト完了")
	response.Success(w, SimilarTracksResponse{Items: similarTracks})
}

// searchSpotifyParallel searches Spotify for tracks by ISRC in parallel.
// It uses a semaphore to limit concurrent requests to maxConcurrent.
// Returns a slice of SimilarTrackResult for successful searches.
func (h *Handler) searchSpotifyParallel(ctx context.Context, isrcList []string) []SimilarTrackResult {
	var results []SimilarTrackResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Semaphore for concurrent limit
	sem := make(chan struct{}, maxConcurrent)

	for _, isrc := range isrcList {
		select {
		case <-ctx.Done():
			return results
		default:
		}

		wg.Add(1)
		go func(isrc string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			// Create context with request timeout
			reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			defer cancel()

			searchResult, err := h.SpotifyClient.SearchByISRC(reqCtx, isrc)
			if err != nil {
				return
			}

			if len(searchResult.Tracks.Items) == 0 {
				return
			}

			track := searchResult.Tracks.Items[0]
			result := convertToSimilarTrackResult(track)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(isrc)
	}

	wg.Wait()
	return results
}

// convertToSimilarTrackResult converts a Spotify SearchTrack to SimilarTrackResult.
func convertToSimilarTrackResult(track spotify.SearchTrack) SimilarTrackResult {
	result := SimilarTrackResult{
		Album: SimilarAlbumResult{
			URL:         track.Album.ExternalURLs["spotify"],
			ID:          track.Album.ID,
			Images:      convertSimilarImages(track.Album.Images),
			Name:        track.Album.Name,
			ReleaseDate: track.Album.ReleaseDate,
			Artists:     convertSimilarAlbumArtists(track.Album.Artists),
		},
		URL:         track.ExternalURLs["spotify"],
		ID:          track.ID,
		Name:        track.Name,
		TrackNumber: track.TrackNumber,
	}

	// Handle nullable fields
	if track.Popularity > 0 {
		result.Popularity = &track.Popularity
	}
	if track.ExternalIDs.ISRC != "" {
		result.ISRC = &track.ExternalIDs.ISRC
	}

	return result
}

// convertSimilarImages converts spotify Image slice to ImageResult slice.
func convertSimilarImages(images []spotify.Image) []ImageResult {
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

// convertSimilarAlbumArtists converts spotify SimpleArtist slice to SimilarArtistResult slice.
func convertSimilarAlbumArtists(artists []spotify.SimpleArtist) []SimilarArtistResult {
	result := make([]SimilarArtistResult, len(artists))
	for i, a := range artists {
		result[i] = SimilarArtistResult{
			URL:  a.ExternalURLs["spotify"],
			Name: a.Name,
			ID:   a.ID,
		}
	}
	return result
}

// removeDuplicates removes duplicate tracks from the slice.
// It uses ISRC for deduplication if available, otherwise uses track ID.
func removeDuplicates(tracks []SimilarTrackResult) []SimilarTrackResult {
	seen := make(map[string]bool)
	result := make([]SimilarTrackResult, 0, len(tracks))

	for _, track := range tracks {
		// Use ISRC for deduplication if available
		key := track.ID
		if track.ISRC != nil && *track.ISRC != "" {
			key = *track.ISRC
		}

		if !seen[key] {
			seen[key] = true
			result = append(result, track)
		}
	}

	return result
}