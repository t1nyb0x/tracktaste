// Package lastfm provides a gateway to the Last.fm API.
package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	baseURL        = "https://ws.audioscrobbler.com/2.0/"
	defaultTimeout = 10 * time.Second
)

// Gateway implements the LastFMAPI interface.
type Gateway struct {
	apiKey     string
	httpClient *http.Client
}

// NewGateway creates a new Last.fm API gateway.
func NewGateway(apiKey string) *Gateway {
	return &Gateway{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// similarTracksResponse represents the Last.fm API response for track.getSimilar.
type similarTracksResponse struct {
	SimilarTracks struct {
		Track []struct {
			Name       string `json:"name"`
			PlayCount  int    `json:"playcount"`
			MBID       string `json:"mbid"`
			Match      string `json:"match"`
			URL        string `json:"url"`
			Streamable struct {
				Text      string `json:"#text"`
				FullTrack string `json:"fulltrack"`
			} `json:"streamable"`
			Duration int `json:"duration"`
			Artist   struct {
				Name string `json:"name"`
				MBID string `json:"mbid"`
				URL  string `json:"url"`
			} `json:"artist"`
		} `json:"track"`
		Attr struct {
			Artist string `json:"artist"`
		} `json:"@attr"`
	} `json:"similartracks"`
	Error   int    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetSimilarTracks returns tracks similar to the given track.
func (g *Gateway) GetSimilarTracks(ctx context.Context, artist, track string, limit int) ([]domain.LastFMTrack, error) {
	params := url.Values{}
	params.Set("method", "track.getSimilar")
	params.Set("artist", artist)
	params.Set("track", track)
	params.Set("api_key", g.apiKey)
	params.Set("format", "json")
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	params.Set("autocorrect", "1") // Auto-correct artist/track names

	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result similarTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API error
	if result.Error != 0 {
		if result.Error == 6 { // Track not found
			logger.Debug("LastFM", fmt.Sprintf("Track not found: %s - %s", artist, track))
			return []domain.LastFMTrack{}, nil
		}
		return nil, fmt.Errorf("Last.fm API error %d: %s", result.Error, result.Message)
	}

	return g.convertToDomain(result.SimilarTracks.Track), nil
}

// GetSimilarTracksByMBID returns tracks similar to the given track using MusicBrainz ID.
func (g *Gateway) GetSimilarTracksByMBID(ctx context.Context, mbid string, limit int) ([]domain.LastFMTrack, error) {
	params := url.Values{}
	params.Set("method", "track.getSimilar")
	params.Set("mbid", mbid)
	params.Set("api_key", g.apiKey)
	params.Set("format", "json")
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result similarTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API error
	if result.Error != 0 {
		if result.Error == 6 { // Track not found
			logger.Debug("LastFM", fmt.Sprintf("Track not found by MBID: %s", mbid))
			return []domain.LastFMTrack{}, nil
		}
		return nil, fmt.Errorf("Last.fm API error %d: %s", result.Error, result.Message)
	}

	return g.convertToDomain(result.SimilarTracks.Track), nil
}

func (g *Gateway) convertToDomain(tracks []struct {
	Name       string `json:"name"`
	PlayCount  int    `json:"playcount"`
	MBID       string `json:"mbid"`
	Match      string `json:"match"`
	URL        string `json:"url"`
	Streamable struct {
		Text      string `json:"#text"`
		FullTrack string `json:"fulltrack"`
	} `json:"streamable"`
	Duration int `json:"duration"`
	Artist   struct {
		Name string `json:"name"`
		MBID string `json:"mbid"`
		URL  string `json:"url"`
	} `json:"artist"`
}) []domain.LastFMTrack {
	result := make([]domain.LastFMTrack, 0, len(tracks))
	for _, t := range tracks {
		match, _ := strconv.ParseFloat(t.Match, 64)
		result = append(result, domain.LastFMTrack{
			Name:       t.Name,
			Artist:     t.Artist.Name,
			MBID:       t.MBID,
			ArtistMBID: t.Artist.MBID,
			Match:      match,
			URL:        t.URL,
		})
	}
	return result
}
