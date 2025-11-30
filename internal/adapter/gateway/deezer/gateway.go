// Package deezer provides the Deezer API gateway implementation.
package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	apiBaseURL = "https://api.deezer.com"
	// Deezer rate limit: 50 requests per 5 seconds
	// We use a conservative limit to avoid hitting rate limits
	maxConcurrentRequests = 10
)

// Gateway implements the DeezerAPI interface.
type Gateway struct {
	httpc *http.Client
}

// NewGateway creates a new Deezer API gateway.
// Deezer API does not require authentication for basic track information.
func NewGateway() *Gateway {
	return &Gateway{
		httpc: &http.Client{Timeout: 10 * time.Second},
	}
}

// rawTrack represents the raw JSON response from Deezer API.
type rawTrack struct {
	ID             int64   `json:"id"`
	Title          string  `json:"title"`
	ISRC           string  `json:"isrc"`
	Duration       int     `json:"duration"` // in seconds
	BPM            float64 `json:"bpm"`
	Gain           float64 `json:"gain"`
	ExplicitLyrics bool    `json:"explicit_lyrics"`
	Artist         *struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// rawSearchResponse represents the search response from Deezer API.
type rawSearchResponse struct {
	Data  []rawTrack `json:"data"`
	Total int        `json:"total"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// GetTrackByISRC searches for a track by ISRC.
func (g *Gateway) GetTrackByISRC(ctx context.Context, isrc string) (*domain.DeezerTrack, error) {
	if isrc == "" {
		return nil, fmt.Errorf("deezer: ISRC is required")
	}

	// Deezer supports ISRC lookup via /track/isrc:{isrc}
	endpoint := fmt.Sprintf("%s/track/isrc:%s", apiBaseURL, isrc)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("deezer: failed to create request: %w", err)
	}

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deezer: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer: unexpected status %d", resp.StatusCode)
	}

	var raw rawTrack
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("deezer: failed to decode response: %w", err)
	}

	// Check for error in response body (Deezer returns 200 with error in body)
	if raw.Error != nil {
		if raw.Error.Code == 800 { // "no data" error
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("deezer: API error: %s", raw.Error.Message)
	}

	return g.convertToTrack(&raw), nil
}

// SearchTrack searches for a track by title and artist.
func (g *Gateway) SearchTrack(ctx context.Context, title, artist string) (*domain.DeezerTrack, error) {
	if title == "" {
		return nil, fmt.Errorf("deezer: title is required")
	}

	// Build search query
	query := fmt.Sprintf("track:\"%s\"", title)
	if artist != "" {
		query += fmt.Sprintf(" artist:\"%s\"", artist)
	}

	endpoint := fmt.Sprintf("%s/search/track?q=%s", apiBaseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("deezer: failed to create request: %w", err)
	}

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deezer: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deezer: unexpected status %d", resp.StatusCode)
	}

	var raw rawSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("deezer: failed to decode response: %w", err)
	}

	if raw.Error != nil {
		return nil, fmt.Errorf("deezer: API error: %s", raw.Error.Message)
	}

	if len(raw.Data) == 0 {
		return nil, domain.ErrNotFound
	}

	// Return the first (most relevant) result
	return g.convertToTrack(&raw.Data[0]), nil
}

// GetTracksByISRCBatch retrieves multiple tracks by their ISRCs.
// Deezer does not have a batch endpoint, so we make parallel requests.
func (g *Gateway) GetTracksByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.DeezerTrack, error) {
	if len(isrcs) == 0 {
		return make(map[string]*domain.DeezerTrack), nil
	}

	result := make(map[string]*domain.DeezerTrack)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrent requests
	sem := make(chan struct{}, maxConcurrentRequests)

	for _, isrc := range isrcs {
		wg.Add(1)
		go func(isrc string) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			track, err := g.GetTrackByISRC(ctx, isrc)
			if err != nil {
				if err != domain.ErrNotFound {
					logger.Warning("Deezer", fmt.Sprintf("Failed to get track by ISRC %s: %v", isrc, err))
				}
				return
			}

			mu.Lock()
			result[isrc] = track
			mu.Unlock()
		}(isrc)
	}

	wg.Wait()

	return result, nil
}

// convertToTrack converts raw Deezer API response to domain model.
func (g *Gateway) convertToTrack(raw *rawTrack) *domain.DeezerTrack {
	track := &domain.DeezerTrack{
		ID:              raw.ID,
		Title:           raw.Title,
		ISRC:            raw.ISRC,
		BPM:             raw.BPM,
		DurationSeconds: raw.Duration,
		Gain:            raw.Gain,
		ExplicitLyrics:  raw.ExplicitLyrics,
	}

	if raw.Artist != nil {
		track.ArtistID = raw.Artist.ID
		track.ArtistName = raw.Artist.Name
	}

	return track
}
