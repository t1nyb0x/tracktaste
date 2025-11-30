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
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/port/repository"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	tokenEndpoint = "https://accounts.spotify.com/api/token"
	apiBaseURL    = "https://api.spotify.com/v1"
)

// isAuthError checks if the status code indicates an authentication error.
// 401 Unauthorized: token is invalid or expired
// 400 Bad Request: can occur with malformed tokens ("Only valid bearer authentication supported")
func isAuthError(statusCode int) bool {
	return statusCode == http.StatusUnauthorized || statusCode == http.StatusBadRequest
}

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

// invalidateToken removes the cached token when API returns an auth error.
func (g *Gateway) invalidateToken(ctx context.Context) {
	if g.tokenRepo != nil {
		if err := g.tokenRepo.InvalidateToken(ctx, "spotify"); err != nil {
			logger.Warning("Spotify", fmt.Sprintf("Failed to invalidate token: %v", err))
		} else {
			logger.Info("Spotify", "Token invalidated due to auth error, will fetch new token on next request")
		}
	}
}

func (g *Gateway) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	return g.getTrackByIDWithRetry(ctx, id, false)
}

func (g *Gateway) getTrackByIDWithRetry(ctx context.Context, id string, isRetry bool) (*domain.Track, error) {
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

	// Retry once with fresh token if auth error
	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getTrackByIDWithRetry(ctx, id, true)
	}

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
	return g.getArtistByIDWithRetry(ctx, id, false)
}

func (g *Gateway) getArtistByIDWithRetry(ctx context.Context, id string, isRetry bool) (*domain.Artist, error) {
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

	// Retry once with fresh token if auth error
	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getArtistByIDWithRetry(ctx, id, true)
	}

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
	return g.getAlbumByIDWithRetry(ctx, id, false)
}

func (g *Gateway) getAlbumByIDWithRetry(ctx context.Context, id string, isRetry bool) (*domain.Album, error) {
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

	// Retry once with fresh token if auth error
	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getAlbumByIDWithRetry(ctx, id, true)
	}

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
	return g.searchTracksWithRetry(ctx, query, false)
}

func (g *Gateway) searchTracksWithRetry(ctx context.Context, query string, isRetry bool) ([]domain.Track, error) {
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

	// Retry once with fresh token if auth error
	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.searchTracksWithRetry(ctx, query, true)
	}

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
	return g.searchByISRCWithRetry(ctx, isrc, false)
}

func (g *Gateway) searchByISRCWithRetry(ctx context.Context, isrc string, isRetry bool) (*domain.Track, error) {
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

	// Retry once with fresh token if auth error
	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.searchByISRCWithRetry(ctx, isrc, true)
	}

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

// GetAudioFeatures retrieves audio features for a single track.
func (g *Gateway) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	return g.getAudioFeaturesWithRetry(ctx, trackID, false)
}

func (g *Gateway) getAudioFeaturesWithRetry(ctx context.Context, trackID string, isRetry bool) (*domain.AudioFeatures, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiBaseURL+"/audio-features/"+trackID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getAudioFeaturesWithRetry(ctx, trackID, true)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify audio-features: status %d", res.StatusCode)
	}

	var raw rawAudioFeatures
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.toDomain(), nil
}

// GetAudioFeaturesBatch retrieves audio features for multiple tracks (max 100).
func (g *Gateway) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	return g.getAudioFeaturesBatchWithRetry(ctx, trackIDs, false)
}

func (g *Gateway) getAudioFeaturesBatchWithRetry(ctx context.Context, trackIDs []string, isRetry bool) ([]domain.AudioFeatures, error) {
	if len(trackIDs) == 0 {
		return []domain.AudioFeatures{}, nil
	}
	if len(trackIDs) > 100 {
		trackIDs = trackIDs[:100]
	}

	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	ids := strings.Join(trackIDs, ",")
	reqURL := fmt.Sprintf("%s/audio-features?ids=%s", apiBaseURL, url.QueryEscape(ids))
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getAudioFeaturesBatchWithRetry(ctx, trackIDs, true)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify audio-features batch: status %d", res.StatusCode)
	}

	var result struct {
		AudioFeatures []*rawAudioFeatures `json:"audio_features"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	features := make([]domain.AudioFeatures, 0, len(result.AudioFeatures))
	for _, raw := range result.AudioFeatures {
		if raw != nil {
			features = append(features, *raw.toDomain())
		}
	}

	return features, nil
}

// GetRecommendations retrieves track recommendations based on seed tracks/artists/genres.
func (g *Gateway) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	return g.getRecommendationsWithRetry(ctx, params, false)
}

func (g *Gateway) getRecommendationsWithRetry(ctx context.Context, params external.RecommendationParams, isRetry bool) ([]domain.Track, error) {
	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	// Build query parameters
	q := url.Values{}
	if len(params.SeedTracks) > 0 {
		q.Set("seed_tracks", strings.Join(params.SeedTracks, ","))
	}
	if len(params.SeedArtists) > 0 {
		q.Set("seed_artists", strings.Join(params.SeedArtists, ","))
	}
	if len(params.SeedGenres) > 0 {
		q.Set("seed_genres", strings.Join(params.SeedGenres, ","))
	}
	if params.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", params.Limit))
	} else {
		q.Set("limit", "20")
	}

	// Target parameters
	if params.TargetTempo != nil {
		q.Set("target_tempo", fmt.Sprintf("%.2f", *params.TargetTempo))
	}
	if params.TargetEnergy != nil {
		q.Set("target_energy", fmt.Sprintf("%.2f", *params.TargetEnergy))
	}
	if params.TargetValence != nil {
		q.Set("target_valence", fmt.Sprintf("%.2f", *params.TargetValence))
	}
	if params.TargetDanceability != nil {
		q.Set("target_danceability", fmt.Sprintf("%.2f", *params.TargetDanceability))
	}
	if params.TargetAcousticness != nil {
		q.Set("target_acousticness", fmt.Sprintf("%.2f", *params.TargetAcousticness))
	}

	// Min/Max parameters
	if params.MinTempo != nil {
		q.Set("min_tempo", fmt.Sprintf("%.2f", *params.MinTempo))
	}
	if params.MaxTempo != nil {
		q.Set("max_tempo", fmt.Sprintf("%.2f", *params.MaxTempo))
	}
	if params.MinEnergy != nil {
		q.Set("min_energy", fmt.Sprintf("%.2f", *params.MinEnergy))
	}
	if params.MaxEnergy != nil {
		q.Set("max_energy", fmt.Sprintf("%.2f", *params.MaxEnergy))
	}

	reqURL := fmt.Sprintf("%s/recommendations?%s", apiBaseURL, q.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getRecommendationsWithRetry(ctx, params, true)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify recommendations: status %d", res.StatusCode)
	}

	var result struct {
		Tracks []rawTrack `json:"tracks"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	tracks := make([]domain.Track, len(result.Tracks))
	for i, raw := range result.Tracks {
		tracks[i] = *raw.toDomain()
	}

	return tracks, nil
}

// GetArtistGenres retrieves genres for a single artist.
func (g *Gateway) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	artist, err := g.GetArtistByID(ctx, artistID)
	if err != nil {
		return nil, err
	}
	return artist.Genres, nil
}

// GetArtistGenresBatch retrieves genres for multiple artists (max 50).
func (g *Gateway) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	return g.getArtistGenresBatchWithRetry(ctx, artistIDs, false)
}

func (g *Gateway) getArtistGenresBatchWithRetry(ctx context.Context, artistIDs []string, isRetry bool) (map[string][]string, error) {
	if len(artistIDs) == 0 {
		return map[string][]string{}, nil
	}
	if len(artistIDs) > 50 {
		artistIDs = artistIDs[:50]
	}

	token, err := g.getToken(ctx)
	if err != nil {
		return nil, err
	}

	ids := strings.Join(artistIDs, ",")
	reqURL := fmt.Sprintf("%s/artists?ids=%s", apiBaseURL, url.QueryEscape(ids))
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := g.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if isAuthError(res.StatusCode) && !isRetry {
		g.invalidateToken(ctx)
		return g.getArtistGenresBatchWithRetry(ctx, artistIDs, true)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify artists batch: status %d", res.StatusCode)
	}

	var result struct {
		Artists []rawArtist `json:"artists"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	genreMap := make(map[string][]string, len(result.Artists))
	for _, artist := range result.Artists {
		genreMap[artist.ID] = artist.Genres
	}

	return genreMap, nil
}
