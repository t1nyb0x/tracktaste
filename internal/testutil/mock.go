// Package testutil provides common test utilities and mocks.
package testutil

import (
	"context"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

// MockSpotifyAPI is a mock implementation of external.SpotifyAPI.
type MockSpotifyAPI struct {
	GetTrackByIDFunc          func(ctx context.Context, id string) (*domain.Track, error)
	GetArtistByIDFunc         func(ctx context.Context, id string) (*domain.Artist, error)
	GetAlbumByIDFunc          func(ctx context.Context, id string) (*domain.Album, error)
	SearchTracksFunc          func(ctx context.Context, query string) ([]domain.Track, error)
	SearchByISRCFunc          func(ctx context.Context, isrc string) (*domain.Track, error)
	GetAudioFeaturesFunc      func(ctx context.Context, trackID string) (*domain.AudioFeatures, error)
	GetAudioFeaturesBatchFunc func(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error)
	GetRecommendationsFunc    func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error)
	GetArtistGenresFunc       func(ctx context.Context, artistID string) ([]string, error)
	GetArtistGenresBatchFunc  func(ctx context.Context, artistIDs []string) (map[string][]string, error)
}

func (m *MockSpotifyAPI) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	if m.GetTrackByIDFunc != nil {
		return m.GetTrackByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	if m.GetArtistByIDFunc != nil {
		return m.GetArtistByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	if m.GetAlbumByIDFunc != nil {
		return m.GetAlbumByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	if m.SearchTracksFunc != nil {
		return m.SearchTracksFunc(ctx, query)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	if m.SearchByISRCFunc != nil {
		return m.SearchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	if m.GetAudioFeaturesFunc != nil {
		return m.GetAudioFeaturesFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	if m.GetAudioFeaturesBatchFunc != nil {
		return m.GetAudioFeaturesBatchFunc(ctx, trackIDs)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	if m.GetRecommendationsFunc != nil {
		return m.GetRecommendationsFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	if m.GetArtistGenresFunc != nil {
		return m.GetArtistGenresFunc(ctx, artistID)
	}
	return nil, nil
}

func (m *MockSpotifyAPI) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	if m.GetArtistGenresBatchFunc != nil {
		return m.GetArtistGenresBatchFunc(ctx, artistIDs)
	}
	return nil, nil
}

// MockKKBOXAPI is a mock implementation of external.KKBOXAPI.
type MockKKBOXAPI struct {
	SearchByISRCFunc         func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error)
	GetRecommendedTracksFunc func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error)
	GetTrackDetailFunc       func(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error)
}

func (m *MockKKBOXAPI) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	if m.SearchByISRCFunc != nil {
		return m.SearchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *MockKKBOXAPI) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	if m.GetRecommendedTracksFunc != nil {
		return m.GetRecommendedTracksFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *MockKKBOXAPI) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	if m.GetTrackDetailFunc != nil {
		return m.GetTrackDetailFunc(ctx, trackID)
	}
	return nil, nil
}

// MockTokenRepository is a mock implementation of repository.TokenRepository.
type MockTokenRepository struct {
	SaveTokenFunc       func(ctx context.Context, key string, token string, ttlSeconds int) error
	GetTokenFunc        func(ctx context.Context, key string) (string, error)
	IsTokenValidFunc    func(ctx context.Context, key string) bool
	InvalidateTokenFunc func(ctx context.Context, key string) error
}

func (m *MockTokenRepository) SaveToken(ctx context.Context, key string, token string, ttlSeconds int) error {
	if m.SaveTokenFunc != nil {
		return m.SaveTokenFunc(ctx, key, token, ttlSeconds)
	}
	return nil
}

func (m *MockTokenRepository) GetToken(ctx context.Context, key string) (string, error) {
	if m.GetTokenFunc != nil {
		return m.GetTokenFunc(ctx, key)
	}
	return "", nil
}

func (m *MockTokenRepository) IsTokenValid(ctx context.Context, key string) bool {
	if m.IsTokenValidFunc != nil {
		return m.IsTokenValidFunc(ctx, key)
	}
	return false
}

func (m *MockTokenRepository) InvalidateToken(ctx context.Context, key string) error {
	if m.InvalidateTokenFunc != nil {
		return m.InvalidateTokenFunc(ctx, key)
	}
	return nil
}

// Helper functions for creating test data

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}

// CreateTestTrack creates a test Track with the given ID and name.
func CreateTestTrack(id, name string) *domain.Track {
	return &domain.Track{
		ID:   id,
		Name: name,
		URL:  "https://open.spotify.com/track/" + id,
		Album: domain.Album{
			ID:   "album1",
			Name: "Test Album",
		},
	}
}

// CreateTestArtist creates a test Artist with the given ID and name.
func CreateTestArtist(id, name string) *domain.Artist {
	return &domain.Artist{
		ID:   id,
		Name: name,
		URL:  "https://open.spotify.com/artist/" + id,
	}
}

// CreateTestAlbum creates a test Album with the given ID and name.
func CreateTestAlbum(id, name string) *domain.Album {
	return &domain.Album{
		ID:   id,
		Name: name,
		URL:  "https://open.spotify.com/album/" + id,
	}
}
