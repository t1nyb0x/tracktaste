package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/repository"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	tokenEndpoint = "https://accounts.spotify.com/api/token"
	apiBaseURL    = "https://api.spotify.com/v1"
)

type Gateway struct {
	clientID  string
	secret    string
	httpc     *http.Client
	tokenRepo repository.TokenRepository
}

func NewGateway(clientID, secret string, tokenRepo repository.TokenRepository) *Gateway {
	return &Gateway{
		clientID:  clientID,
		secret:    secret,
		httpc:     &http.Client{Timeout: 10 * time.Second},
		tokenRepo: tokenRepo,
	}
}

func (g *Gateway) getToken(ctx context.Context) (string, error) {
	if g.tokenRepo != nil && g.tokenRepo.IsTokenValid(ctx, "spotify") {
		token, err := g.tokenRepo.GetToken(ctx, "spotify")
		if err == nil {
			return token, nil
		}
	}

	token, expiresIn, err := g.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	if g.tokenRepo != nil {
		if err := g.tokenRepo.SaveToken(ctx, "spotify", token, expiresIn); err != nil {
			logger.Warning("Spotify", fmt.Sprintf("Failed to save token: %v", err))
		}
	}

	return token, nil
}

func (g *Gateway) fetchToken(ctx context.Context) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", g.clientID)
	data.Set("client_secret", g.secret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := g.httpc.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("spotify token: status %d", res.StatusCode)
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", 0, err
	}
	return resp.AccessToken, resp.ExpiresIn, nil
}

func (g *Gateway) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var raw rawTrack
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.toDomain(), nil
}

func (g *Gateway) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/artists/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var raw rawArtist
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.toDomain(), nil
}

func (g *Gateway) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/albums/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var raw rawAlbum
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.toDomain(), nil
}

func (g *Gateway) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/search?q=%s&type=track&limit=20", apiBaseURL, url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify search: status %d", res.StatusCode)
	}

	var result struct {
		Tracks struct {
			Items []rawTrack `json:"items"`
		} `json:"tracks"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	tracks := make([]domain.Track, len(result.Tracks.Items))
	for i, raw := range result.Tracks.Items {
		tracks[i] = *raw.toDomain()
	}
	return tracks, nil
}

func (g *Gateway) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/search?q=isrc:%s&type=track&limit=1", apiBaseURL, url.QueryEscape(isrc))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify search: status %d", res.StatusCode)
	}

	var result struct {
		Tracks struct {
			Items []rawTrack `json:"items"`
		} `json:"tracks"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Tracks.Items) == 0 {
		return nil, nil
	}

	return result.Tracks.Items[0].toDomain(), nil
}
