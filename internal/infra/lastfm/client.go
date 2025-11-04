package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey string
	httpc  *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpc:  &http.Client{Timeout: 7 * time.Second},
	}
}

func (c *Client) get(ctx context.Context, params url.Values, v any) error {
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ws.audioscrobbler.com/2.0/?"+params.Encode(), nil)
	if err != nil { return err }
	res, err := c.httpc.Do(req)
	if err != nil { return err }
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("lastfm: status %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(v)
}

// ---- mapping
type img struct {
	Size string `json:"size"`
	URL  string `json:"#text"`
}

func pickImages(imgs []img) (small, large string) {
	for _, im := range imgs {
		if im.Size == "small" && im.URL != "" { small = im.URL }
		if (im.Size == "extralarge" || im.Size == "mega") && im.URL != "" { large = im.URL }
	}
	return
}
