package kkbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/port/repository"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	tokenEndpoint = "https://account.kkbox.com/oauth2/token"
	apiBaseURL    = "https://api.kkbox.com/v1.1"
	territory     = "JP"
)

type Gateway struct {
	clientID     string
	clientSecret string
	httpc        *http.Client
	tokenRepo    repository.TokenRepository
}

func NewGateway(clientID, clientSecret string, tokenRepo repository.TokenRepository) *Gateway {
	return &Gateway{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpc:        &http.Client{Timeout: 10 * time.Second},
		tokenRepo:    tokenRepo,
	}
}

func (g *Gateway) getToken(ctx context.Context) (string, error) {
	if g.tokenRepo != nil && g.tokenRepo.IsTokenValid(ctx, "kkbox") {
		token, err := g.tokenRepo.GetToken(ctx, "kkbox")
		if err == nil {
			return token, nil
		}
	}

	token, expiresIn, err := g.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	if g.tokenRepo != nil {
		if err := g.tokenRepo.SaveToken(ctx, "kkbox", token, expiresIn); err != nil {
			logger.Warning("KKBOX", fmt.Sprintf("Failed to save token: %v", err))
		}
	}

	return token, nil
}

func (g *Gateway) fetchToken(ctx context.Context) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", g.clientID)
	data.Set("client_secret", g.clientSecret)

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
		return "", 0, fmt.Errorf("kkbox token: status %d", res.StatusCode)
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

func (g *Gateway) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/search?q=%s&type=track&territory=%s&limit=1", apiBaseURL, url.QueryEscape(isrc), territory)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
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
		return nil, fmt.Errorf("kkbox search: status %d", res.StatusCode)
	}

	var result struct {
		Tracks struct {
			Data []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				ISRC string `json:"isrc"`
			} `json:"data"`
		} `json:"tracks"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Tracks.Data) == 0 {
		return nil, nil
	}

	t := result.Tracks.Data[0]
	return &external.KKBOXTrackInfo{ID: t.ID, Name: t.Name, ISRC: t.ISRC}, nil
}

func (g *Gateway) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/tracks/%s/related-tracks?territory=%s&limit=50", apiBaseURL, trackID, territory)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
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
		return nil, fmt.Errorf("kkbox recommend: status %d", res.StatusCode)
	}

	var result struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	tracks := make([]external.KKBOXTrackInfo, len(result.Data))
	for i, t := range result.Data {
		tracks[i] = external.KKBOXTrackInfo{ID: t.ID, Name: t.Name}
	}
	return tracks, nil
}

func (g *Gateway) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/tracks/%s?territory=%s", apiBaseURL, trackID, territory)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
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
		return nil, fmt.Errorf("kkbox detail: status %d", res.StatusCode)
	}

	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		ISRC string `json:"isrc"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &external.KKBOXTrackInfo{ID: result.ID, Name: result.Name, ISRC: result.ISRC}, nil
}
