package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) GetBearerToken(ctx context.Context, params url.Values) (string, error) {
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", c.APIKey)
	params.Set("client_secret", c.Secret)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil { return "", err }
	res, err := c.httpc.Do(req)
	if err != nil { return "", err }
	defer res.Body.Close()
	if res.StatusCode != 200 {
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



func (c *Client) get(ctx context.Context, params url.Values, v any) error {
	params.Set("bearer_token", c.bearerToken)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.spotify.com/v1/tracks/" + params.Encode(), nil)
	if err != nil { return err }
	res, err := c.httpc.Do(req)
	if err != nil { return err }
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("spotify: status %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(v)
}