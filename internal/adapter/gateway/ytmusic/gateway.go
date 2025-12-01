// Package ytmusic provides a gateway to interact with the YouTube Music sidecar service.
package ytmusic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const featureName = "YTMusic"

// Gateway implements the YouTubeMusicAPI interface using the Python sidecar service.
type Gateway struct {
	baseURL    string
	httpClient *http.Client
}

// NewGateway creates a new YouTube Music gateway.
// baseURL should be the sidecar service URL (e.g., "http://ytmusic-sidecar:8081")
func NewGateway(baseURL string) *Gateway {
	return &Gateway{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// similarResponse represents the JSON response from /similar endpoint.
type similarResponse struct {
	VideoID string      `json:"video_id"`
	Tracks  []trackJSON `json:"tracks"`
}

// searchResponse represents the JSON response from /search endpoint.
type searchResponse struct {
	Query  string      `json:"query"`
	Tracks []trackJSON `json:"tracks"`
}

// trackJSON represents a track in the sidecar API response.
type trackJSON struct {
	VideoID         string  `json:"video_id"`
	Title           string  `json:"title"`
	Artist          string  `json:"artist"`
	ArtistID        *string `json:"artist_id"`
	Album           *string `json:"album"`
	AlbumID         *string `json:"album_id"`
	DurationSeconds *int    `json:"duration_seconds"`
	ThumbnailURL    *string `json:"thumbnail_url"`
	IsExplicit      bool    `json:"is_explicit"`
}

// GetSimilarTracks retrieves similar tracks for a given YouTube video ID.
func (g *Gateway) GetSimilarTracks(ctx context.Context, videoID string, limit int) ([]domain.YTMusicTrack, error) {
	logger.Debug(featureName, fmt.Sprintf("getting similar tracks for videoID=%s, limit=%d", videoID, limit))

	reqURL := fmt.Sprintf("%s/similar/%s?limit=%d", g.baseURL, url.PathEscape(videoID), limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar tracks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sidecar returned status %d", resp.StatusCode)
	}

	var result similarResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	tracks := make([]domain.YTMusicTrack, 0, len(result.Tracks))
	for _, t := range result.Tracks {
		tracks = append(tracks, convertTrack(t))
	}

	logger.Debug(featureName, fmt.Sprintf("found %d similar tracks", len(tracks)))
	return tracks, nil
}

// SearchTracks searches for tracks on YouTube Music.
func (g *Gateway) SearchTracks(ctx context.Context, query string, limit int) ([]domain.YTMusicTrack, error) {
	logger.Debug(featureName, fmt.Sprintf("searching for query=%s, limit=%d", query, limit))

	reqURL := fmt.Sprintf("%s/search?q=%s&limit=%d", g.baseURL, url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sidecar returned status %d", resp.StatusCode)
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	tracks := make([]domain.YTMusicTrack, 0, len(result.Tracks))
	for _, t := range result.Tracks {
		tracks = append(tracks, convertTrack(t))
	}

	logger.Debug(featureName, fmt.Sprintf("found %d tracks for query", len(tracks)))
	return tracks, nil
}

// convertTrack converts a trackJSON to domain.YTMusicTrack.
func convertTrack(t trackJSON) domain.YTMusicTrack {
	track := domain.YTMusicTrack{
		VideoID:    t.VideoID,
		Title:      t.Title,
		Artist:     t.Artist,
		IsExplicit: t.IsExplicit,
	}

	if t.ArtistID != nil {
		track.ArtistID = *t.ArtistID
	}
	if t.Album != nil {
		track.Album = *t.Album
	}
	if t.AlbumID != nil {
		track.AlbumID = *t.AlbumID
	}
	if t.DurationSeconds != nil {
		track.DurationSeconds = *t.DurationSeconds
	}
	if t.ThumbnailURL != nil {
		track.ThumbnailURL = *t.ThumbnailURL
	}

	return track
}
