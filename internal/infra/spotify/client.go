package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	APIKey string
	Secret string
	httpc  *http.Client
}

func New(APIKey string, Secret string) *Client {
	return &Client{
		APIKey: APIKey,
		Secret: Secret,
		httpc: &http.Client{Timeout: 7 * time.Second},
	}
}

func (c *Client) GetBearerToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.APIKey)
	data.Set("client_secret", c.Secret)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://accounts.spotify.com/api/token",
		strings.NewReader(data.Encode()),
	)
	
	if err != nil { return "", err }

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpc.Do(req)

	if err != nil { return "", err }
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}



func (c *Client) FetchById(ctx context.Context, params string, v any) error {
	log.Println("[DEBUG] FetchById called with params:", params)
	// Get bearer token with required arguments
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Obtained bearer token")

	// Create request with proper authorization header
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/tracks/" + params, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer " + token)

	res, err := c.httpc.Do(req)
	if err != nil { return err }
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify: status %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}
