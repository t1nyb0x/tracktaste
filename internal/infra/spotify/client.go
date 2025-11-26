// Package spotify provides a client for interacting with the Spotify Web API.
// It supports authentication, track/artist/album retrieval, and search functionality.
// Token caching is handled via Redis for optimal performance.
package spotify

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
	tokenEndpoint = "https://accounts.spotify.com/api/token"
	apiBaseURL    = "https://api.spotify.com/v1"
)

// Client represents a Spotify API client with authentication credentials.
type Client struct {
	APIKey string
	Secret string
	httpc  *http.Client
}

// New creates a new Spotify API client with the provided credentials.
// The client has a default HTTP timeout of 10 seconds.
//
// Parameters:
//   - APIKey: Spotify OAuth client ID
//   - Secret: Spotify OAuth client secret
//
// Returns a configured Spotify client instance.
func New(APIKey string, Secret string) *Client {
	return &Client{
		APIKey: APIKey,
		Secret: Secret,
		httpc:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetBearerToken retrieves a bearer token for Spotify API authentication.
// It first checks Redis cache for an existing valid token. If not found or expired,
// it fetches a new token from Spotify OAuth endpoint and caches it in Redis.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns the bearer token string or an error if token retrieval fails.
func (c *Client) GetBearerToken(ctx context.Context) (string, error) {
	// Try to get from Redis
	if redisClient.IsTokenValid(ctx, "spotify") {
		token, err := redisClient.GetToken(ctx, "spotify")
		if err == nil {
			logger.Debug("Spotify", "Token retrieved from Redis")
			return token, nil
		}
	}

	// Fetch new token
	token, expiresIn, err := c.fetchToken(ctx)
	if err != nil {
		logger.Error("Spotify", fmt.Sprintf("Failed to fetch token: %v", err))
		return "", err
	}

	// Save to Redis
	if err := redisClient.SaveToken(ctx, "spotify", token, expiresIn); err != nil {
		logger.Warning("Spotify", fmt.Sprintf("Failed to save token to Redis: %v", err))
	}

	logger.Info("Spotify", "New token fetched and saved")
	return token, nil
}

func (c *Client) fetchToken(ctx context.Context) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.APIKey)
	data.Set("client_secret", c.Secret)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		tokenEndpoint,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpc.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", 0, err
	}
	return resp.AccessToken, resp.ExpiresIn, nil
}

// FetchById fetches track information by Spotify track ID.
// The result is decoded into the provided interface.
//
// Parameters:
//   - ctx: Context for the operation
//   - params: Spotify track ID
//   - v: Pointer to struct to decode the response into
//
// Returns an error if the request or decoding fails.
func (c *Client) FetchById(ctx context.Context, params string, v any) error {
	logger.Debug("Spotify", fmt.Sprintf("FetchById called with params: %s", params))
	
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return err
	}
	logger.Debug("Spotify", "Obtained bearer token")

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/tracks/"+params, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}

// FetchArtistById fetches artist information by Spotify artist ID.
// Returns detailed artist data including followers, genres, and images.
//
// Parameters:
//   - ctx: Context for the operation
//   - id: Spotify artist ID
//
// Returns artist response or an error if the request fails.
func (c *Client) FetchArtistById(ctx context.Context, id string) (*ArtistResponse, error) {
	logger.Debug("Spotify", fmt.Sprintf("FetchArtistById called with id: %s", id))

	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/artists/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var artist ArtistResponse
	if err := json.NewDecoder(res.Body).Decode(&artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

// FetchAlbumById fetches album information by Spotify album ID.
// Returns detailed album data including tracks, artists, and images.
//
// Parameters:
//   - ctx: Context for the operation
//   - id: Spotify album ID
//
// Returns album response or an error if the request fails.
func (c *Client) FetchAlbumById(ctx context.Context, id string) (*AlbumResponse, error) {
	logger.Debug("Spotify", fmt.Sprintf("FetchAlbumById called with id: %s", id))

	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/albums/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var album AlbumResponse
	if err := json.NewDecoder(res.Body).Decode(&album); err != nil {
		return nil, err
	}

	return &album, nil
}

// SearchByISRC searches for tracks by International Standard Recording Code (ISRC).
// This is useful for finding the same track across different platforms.
//
// Parameters:
//   - ctx: Context for the operation
//   - isrc: ISRC code to search for
//
// Returns search results or an error if the search fails.
func (c *Client) SearchByISRC(ctx context.Context, isrc string) (*SearchResponse, error) {
	return c.search(ctx, fmt.Sprintf("isrc:%s", isrc))
}

// SearchByUPC searches for tracks by Universal Product Code (UPC).
// UPC is typically associated with album releases.
//
// Parameters:
//   - ctx: Context for the operation
//   - upc: UPC code to search for
//
// Returns search results or an error if the search fails.
func (c *Client) SearchByUPC(ctx context.Context, upc string) (*SearchResponse, error) {
	return c.search(ctx, fmt.Sprintf("upc:%s", upc))
}

func (c *Client) search(ctx context.Context, query string) (*SearchResponse, error) {
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	searchURL := fmt.Sprintf("%s/search?q=%s&type=track&limit=1", apiBaseURL, url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
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
		return nil, fmt.Errorf("spotify API error: status %d", res.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ArtistResponse represents the Spotify Artist API response
type ArtistResponse struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Followers    struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres     []string `json:"genres"`
	Href       string   `json:"href"`
	ID         string   `json:"id"`
	Images     []Image  `json:"images"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
	URI        string   `json:"uri"`
}

// AlbumResponse represents the Spotify Album API response
type AlbumResponse struct {
	AlbumType    string            `json:"album_type"`
	TotalTracks  int               `json:"total_tracks"`
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Images       []Image           `json:"images"`
	Name         string            `json:"name"`
	ReleaseDate  string            `json:"release_date"`
	Artists      []SimpleArtist    `json:"artists"`
	Tracks       TracksPage        `json:"tracks"`
	Popularity   int               `json:"popularity"`
	ExternalIDs  struct {
		UPC string `json:"upc"`
	} `json:"external_ids"`
	Genres []string `json:"genres"`
}

// Image represents an image object
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// SimpleArtist represents a simplified artist object
type SimpleArtist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

// TracksPage represents a paginated tracks response
type TracksPage struct {
	Items []SimpleTrack `json:"items"`
}

// SimpleTrack represents a simplified track object
type SimpleTrack struct {
	Artists     []SimpleArtist    `json:"artists"`
	ExternalURLs map[string]string `json:"external_urls"`
	Href        string            `json:"href"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	TrackNumber int               `json:"track_number"`
	URI         string            `json:"uri"`
}

// SearchResponse represents the Spotify Search API response
type SearchResponse struct {
	Tracks struct {
		Items []SearchTrack `json:"items"`
	} `json:"tracks"`
}

// SearchTrack represents a track in search results
type SearchTrack struct {
	Album        SearchAlbum       `json:"album"`
	Artists      []SimpleArtist    `json:"artists"`
	ExternalIDs  struct {
		ISRC string `json:"isrc"`
	} `json:"external_ids"`
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	TrackNumber  int               `json:"track_number"`
}

// SearchAlbum represents an album in search results
type SearchAlbum struct {
	ExternalURLs map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Images       []Image           `json:"images"`
	Name         string            `json:"name"`
	ReleaseDate  string            `json:"release_date"`
	Artists      []SimpleArtist    `json:"artists"`
}
