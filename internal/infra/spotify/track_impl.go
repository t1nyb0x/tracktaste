package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	domain "github.com/t1nyb0x/tracktaste/internal/domain/spotify"
)

func (c *Client) GetInfo(ctx context.Context, id string) (domain.Track, error) {


	var raw domain.Track
	if err := c.FetchById(ctx, id, &raw); err != nil {
		return domain.Track{}, err
	}
	log.Println("[DEBUG] Fetched track info:", raw)

	artists := make([]domain.Artist, len(raw.Artists))
	for i, a := range raw.Artists {
		artists[i] = domain.Artist{
			ArtistName: a.ArtistName,
			ArtistExternalURLs: a.ArtistExternalURLs,
			Href: a.Href,
			ID: a.ID,
			Typee: a.Typee,
			URI: a.URI,
		}
	}

	availableMarkets := make([]string, len(raw.AvailableMarkets))
	copy(availableMarkets, raw.AvailableMarkets)

	linkedFrom := make(map[string]string)
	for k, v := range raw.LinkedFrom {
		linkedFrom[k] = v
	}

	return domain.Track{
		Album: raw.Album,
		Artists:  artists,
		AvailableMarkets: availableMarkets,
		DiscNumber: raw.DiscNumber,
		DurationMs: raw.DurationMs,
		Explicit:   raw.Explicit,
		TrackExternalIDs:        domain.TrackExternalIDs{
			ISRC: raw.TrackExternalIDs.ISRC,
		},
		TrackExternalUrls: domain.TrackExternalUrls{
			SpotifyTrackURL: raw.TrackExternalUrls.SpotifyTrackURL,
		},
		Href:       raw.Href,
		ID:         raw.ID,
		IsPlayable: raw.IsPlayable,
		LinkedFrom: linkedFrom,
		Name:        raw.Name,
		Popularity: raw.Popularity,
		PreviewURL: raw.PreviewURL,
		TrackNumber: raw.TrackNumber,
		Typee:      raw.Typee,
		URI:        raw.URI,
		IsLocal:    raw.IsLocal,
	}, nil
}

func (c *Client) Search(ctx context.Context, query string) ([]domain.Track, error) {
	// Get bearer token
	token, err := c.GetBearerToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %w", err)
	}

	// Build search URL
	u, err := url.Parse("https://api.spotify.com/v1/search")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	q := u.Query()
	q.Set("q", query)
	q.Set("type", "track")
	q.Set("limit", "10")
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute request
	res, err := c.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify search failed with status: %d", res.StatusCode)
	}

	// Parse response
	var response struct {
		Tracks struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Artists []struct {
					Name string `json:"name"`
					ExternalURLs struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
				} `json:"artists"`
				Album struct {
					Name         string `json:"name"`
					ReleaseDate  string `json:"release_date"`
					ExternalURLs struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
				} `json:"album"`
				ExternalIDs struct {
					ISRC string `json:"isrc"`
				} `json:"external_ids"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map to domain objects
	tracks := make([]domain.Track, len(response.Tracks.Items))
	for i, item := range response.Tracks.Items {
		// var artistName, artistURL string
		// if len(item.Artists) > 0 {
		// 	artistName = item.Artists[0].Name
		// 	artistURL = item.Artists[0].ExternalURLs.Spotify
		// }

		tracks[i] = domain.Track{
			Name:        item.Name,
			// TrackURL:    item.ExternalURLs.Spotify,
			// ISRC:        item.ExternalIDs.ISRC,
			// ArtistName:  artistName,
			// ArtistURL:   artistURL,
			// AlbumName:   item.Album.Name,
			// AlbumURL:    item.Album.ExternalURLs.Spotify,
			// ReleaseDate: item.Album.ReleaseDate,
		}
	}

	return tracks, nil
}
