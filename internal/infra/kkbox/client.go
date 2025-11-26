// Package kkbox provides a client for interacting with the KKBOX API.
// It supports authentication, track search, and recommendation features.
// All API requests use territory=JP (Japan) by default.
package kkbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	redisClient "github.com/t1nyb0x/tracktaste/internal/infra/redis"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	tokenEndpoint = "https://account.kkbox.com/oauth2/token"
	apiBaseURL    = "https://api.kkbox.com/v1.1"
	territory     = "JP"
)

// Client represents a KKBOX API client with authentication credentials.
type Client struct {
	ClientID     string
	ClientSecret string
	httpc        *http.Client
}

// TokenResponse represents the OAuth token response from KKBOX.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SearchResponse represents the search API response containing tracks.
type SearchResponse struct {
	Tracks struct {
		Data []KKBOXTrack `json:"data"`
	} `json:"tracks"`
}

// KKBOXTrack represents a track object from KKBOX API.
type KKBOXTrack struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	URL            string     `json:"url"`
	Album          KKBOXAlbum `json:"album"`
	ISRC           string     `json:"isrc"`
	ExplicitLyrics bool       `json:"explicit_lyrics"`
}

// KKBOXAlbum represents an album object from KKBOX API.
type KKBOXAlbum struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RecommendResponse represents the recommended tracks API response.
type RecommendResponse struct {
	Data []KKBOXTrack `json:"data"`
}

// TrackDetailResponse represents the detailed track information response.
type TrackDetailResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	ISRC  string `json:"isrc"`
	Album struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"album"`
}

// New creates a new KKBOX API client with the provided credentials.
// The client has a default HTTP timeout of 10 seconds.
//
// Parameters:
//   - clientID: KKBOX OAuth client ID
//   - clientSecret: KKBOX OAuth client secret
//
// Returns a configured KKBOX client instance.
func New(clientID string, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		httpc:        &http.Client{Timeout: 10 * time.Second},
	}
}

// GetBearerToken retrieves a bearer token for KKBOX API authentication.
// It first checks Redis cache for an existing valid token. If not found or expired,
// it fetches a new token from KKBOX OAuth endpoint and caches it in Redis.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns the bearer token string or an error if token retrieval fails.
func (c *Client) GetBearerToken(ctx context.Context) (string, error) {
	// Try to get from Redis
	if redisClient.IsTokenValid(ctx, "kkbox") {
		token, err := redisClient.GetToken(ctx, "kkbox")
		if err == nil {
			logger.Debug("KKBOX", "Token retrieved from Redis")
			return token, nil
		}
	}

	// Fetch new token
	token, expiresIn, err := c.fetchToken(ctx)
	if err != nil {
		logger.Error("KKBOX", fmt.Sprintf("Failed to fetch token: %v", err))
		return "", err
	}

	// Save to Redis
	if err := redisClient.SaveToken(ctx, "kkbox", token, expiresIn); err != nil {
		logger.Warning("KKBOX", fmt.Sprintf("Failed to save token to Redis: %v", err))
	}

	logger.Info("KKBOX", "New token fetched and saved")
	return token, nil
}

func (c *Client) fetchToken(ctx context.Context) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpc.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch token: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("failed to fetch token: status %d", res.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}

// Search searches for tracks in KKBOX by the given query string.
// Typically used to search by ISRC or track name.
// Results are limited to 1 track and territory is set to JP.
//
// Parameters:
//   - ctx: Context for the operation
//   - query: Search query (e.g., ISRC code or track name)
//
// Returns search results or an error if the search fails.
func (c *Client) Search(ctx context.Context, query string) (*SearchResponse, error) {
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/search?q=%s&type=track&territory=%s&limit=1",
		apiBaseURL, url.QueryEscape(query), territory)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KKBOX API error: status %d", res.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetRecommendedTracks retrieves recommended tracks for a given KKBOX track ID.
// Returns up to 50 related tracks based on KKBOX's recommendation algorithm.
//
// Parameters:
//   - ctx: Context for the operation
//   - trackID: KKBOX track ID to get recommendations for
//
// Returns recommended tracks or an error if the request fails.
func (c *Client) GetRecommendedTracks(ctx context.Context, trackID string) (*RecommendResponse, error) {
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/tracks/%s/related-tracks?territory=%s&limit=50",
		apiBaseURL, trackID, territory)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KKBOX API error: status %d", res.StatusCode)
	}

	var result RecommendResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetTrackDetail retrieves detailed track information from KKBOX.
// This includes ISRC code which is essential for cross-platform track matching.
//
// Parameters:
//   - ctx: Context for the operation
//   - trackID: KKBOX track ID to get details for
//
// Returns track details or an error if the request fails.
func (c *Client) GetTrackDetail(ctx context.Context, trackID string) (*TrackDetailResponse, error) {
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/tracks/%s?territory=%s", apiBaseURL, trackID, territory)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KKBOX API error: status %d", res.StatusCode)
	}

	var result TrackDetailResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
