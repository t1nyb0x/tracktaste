package lastfm

import (
	"context"
	"net/url"
	"strconv"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/lastfm"
)

// GetInfo
type getInfoResp struct {
	Artist struct {
		Name      string `json:"name"`
		Mbid      string `json:"mbid"`
		URL       string `json:"url"`
		Image     []img  `json:"image"`
		Stats     struct {
			Listeners string `json:"listeners"`
			Playcount string `json:"playcount"`
		} `json:"stats"`
	} `json:"artist"`
}

func (c *Client) GetInfo(ctx context.Context, name string) (domain.Artist, error) {
	params := url.Values{}
	params.Set("method", "artist.getinfo")
	params.Set("artist", name)

	var raw getInfoResp
	if err := c.get(ctx, params, &raw); err != nil {
		return domain.Artist{}, err
	}
	small, large := pickImages(raw.Artist.Image)
	l, _ := strconv.ParseInt(raw.Artist.Stats.Listeners, 10, 64)
	p, _ := strconv.ParseInt(raw.Artist.Stats.Playcount, 10, 64)
	return domain.Artist{
		Name: raw.Artist.Name, Mbid: raw.Artist.Mbid, Url: raw.Artist.URL,
		ImageSmall: small, ImageLarge: large, Listeners: l, Playcount: p,
	}, nil
}

// SearchFirst
type searchResp struct {
	Results struct {
		ArtistMatches struct {
			Artist []struct {
				Name  string `json:"name"`
				Mbid  string `json:"mbid"`
				URL   string `json:"url"`
				Image []img  `json:"image"`
			} `json:"artist"`
		} `json:"artistmatches"`
	} `json:"results"`
}

func (c *Client) SearchFirst(ctx context.Context, query string) (domain.Artist, error) {
	params := url.Values{}
	params.Set("method", "artist.search")
	params.Set("artist", query)
	params.Set("limit", "1")

	var raw searchResp
	if err := c.get(ctx, params, &raw); err != nil {
		return domain.Artist{}, err
	}
	if len(raw.Results.ArtistMatches.Artist) == 0 {
		return domain.Artist{}, nil
	}
	a := raw.Results.ArtistMatches.Artist[0]
	small, large := pickImages(a.Image)
	return domain.Artist{
		Name: a.Name, Mbid: a.Mbid, Url: a.URL,
		ImageSmall: small, ImageLarge: large,
	}, nil
}
