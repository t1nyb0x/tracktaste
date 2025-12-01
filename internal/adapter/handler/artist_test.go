package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	usecasev1 "github.com/t1nyb0x/tracktaste/internal/usecase/v1"
)

// mockSpotifyAPIForArtist for artist handler tests
type mockSpotifyAPIForArtist struct {
	GetArtistByIDFunc func(ctx context.Context, id string) (*domain.Artist, error)
}

func (m *mockSpotifyAPIForArtist) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	if m.GetArtistByIDFunc != nil {
		return m.GetArtistByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForArtist) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	return nil, nil
}

var _ external.SpotifyAPI = (*mockSpotifyAPIForArtist)(nil)

func createTestArtist() *domain.Artist {
	popularity := 85
	followers := 1000000
	return &domain.Artist{
		ID:         "test-artist-id",
		Name:       "Test Artist",
		URL:        "https://open.spotify.com/artist/test-artist-id",
		Popularity: &popularity,
		Followers:  &followers,
		Genres:     []string{"rock", "pop"},
		Images: []domain.Image{
			{URL: "https://example.com/artist.jpg", Height: 640, Width: 640},
		},
	}
}

func TestArtistHandler_FetchByURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockFunc       func(ctx context.Context, id string) (*domain.Artist, error)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "正常系: 有効なURL",
			url:  "https://open.spotify.com/artist/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Artist, error) {
				return createTestArtist(), nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "正常系: intl-ja付きURL",
			url:  "https://open.spotify.com/intl-ja/artist/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Artist, error) {
				return createTestArtist(), nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "異常系: 空のURL",
			url:            "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "EMPTY_PARAM",
		},
		{
			name:           "異常系: Spotify以外のURL",
			url:            "https://music.apple.com/artist/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "NOT_SPOTIFY_URL",
		},
		{
			name:           "異常系: trackのURL",
			url:            "https://open.spotify.com/track/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:           "異常系: albumのURL",
			url:            "https://open.spotify.com/album/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "DIFFERENT_SPOTIFY_URL",
		},
		{
			name: "異常系: APIエラー",
			url:  "https://open.spotify.com/artist/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Artist, error) {
				return nil, errors.New("API error")
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SOMETHING_SPOTIFY_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &mockSpotifyAPIForArtist{GetArtistByIDFunc: tt.mockFunc}
			artistUC := usecasev1.NewArtistUseCase(mockAPI)
			handler := NewArtistHandler(artistUC)

			req := httptest.NewRequest(http.MethodGet, "/v1/artist/fetch?url="+tt.url, nil)
			rec := httptest.NewRecorder()

			handler.FetchByURL(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var resp map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if code, ok := resp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("expected code %s, got %v", tt.expectedCode, resp["code"])
				}
			}

			if tt.expectedStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				result, ok := resp["result"].(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be object")
				}
				if result["id"] != "test-artist-id" {
					t.Errorf("expected id test-artist-id, got %v", result["id"])
				}
				if result["name"] != "Test Artist" {
					t.Errorf("expected name Test Artist, got %v", result["name"])
				}
				// Check followers is formatted as string
				if result["followers"] != "1000000" {
					t.Errorf("expected followers '1000000', got %v", result["followers"])
				}
			}
		})
	}
}

func TestArtistHandler_FetchByURL_NilFollowers(t *testing.T) {
	mockAPI := &mockSpotifyAPIForArtist{
		GetArtistByIDFunc: func(ctx context.Context, id string) (*domain.Artist, error) {
			return &domain.Artist{
				ID:        "test-id",
				Name:      "Test",
				URL:       "https://open.spotify.com/artist/test-id",
				Followers: nil,
				Genres:    []string{},
				Images:    []domain.Image{},
			}, nil
		},
	}
	artistUC := usecasev1.NewArtistUseCase(mockAPI)
	handler := NewArtistHandler(artistUC)

	req := httptest.NewRequest(http.MethodGet, "/v1/artist/fetch?url=https://open.spotify.com/artist/abc123", nil)
	rec := httptest.NewRecorder()

	handler.FetchByURL(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be object")
	}
	// Followers should default to "0" when nil
	if result["followers"] != "0" {
		t.Errorf("expected followers '0' for nil, got %v", result["followers"])
	}
}
