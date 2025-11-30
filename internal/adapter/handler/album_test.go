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
	"github.com/t1nyb0x/tracktaste/internal/usecase"
)

// mockSpotifyAPIForAlbum for album handler tests
type mockSpotifyAPIForAlbum struct {
	GetAlbumByIDFunc func(ctx context.Context, id string) (*domain.Album, error)
}

func (m *mockSpotifyAPIForAlbum) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	if m.GetAlbumByIDFunc != nil {
		return m.GetAlbumByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	return nil, nil
}

func (m *mockSpotifyAPIForAlbum) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	return nil, nil
}

var _ external.SpotifyAPI = (*mockSpotifyAPIForAlbum)(nil)

func createTestAlbum() *domain.Album {
	popularity := 75
	upc := "123456789012"
	return &domain.Album{
		ID:          "test-album-id",
		Name:        "Test Album",
		URL:         "https://open.spotify.com/album/test-album-id",
		ReleaseDate: "2024-01-15",
		Popularity:  &popularity,
		UPC:         &upc,
		Genres:      []string{"rock"},
		Images: []domain.Image{
			{URL: "https://example.com/album.jpg", Height: 640, Width: 640},
		},
		Artists: []domain.Artist{
			{ID: "artist-1", Name: "Test Artist", URL: "https://open.spotify.com/artist/artist-1"},
		},
		Tracks: []domain.SimpleTrack{
			{
				ID:          "track-1",
				Name:        "Track One",
				URL:         "https://open.spotify.com/track/track-1",
				TrackNumber: 1,
				Artists: []domain.Artist{
					{ID: "artist-1", Name: "Test Artist", URL: "https://open.spotify.com/artist/artist-1"},
				},
			},
			{
				ID:          "track-2",
				Name:        "Track Two",
				URL:         "https://open.spotify.com/track/track-2",
				TrackNumber: 2,
				Artists: []domain.Artist{
					{ID: "artist-1", Name: "Test Artist", URL: "https://open.spotify.com/artist/artist-1"},
				},
			},
		},
	}
}

func TestAlbumHandler_FetchByURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockFunc       func(ctx context.Context, id string) (*domain.Album, error)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "正常系: 有効なURL",
			url:  "https://open.spotify.com/album/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Album, error) {
				return createTestAlbum(), nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "正常系: intl-ja付きURL",
			url:  "https://open.spotify.com/intl-ja/album/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Album, error) {
				return createTestAlbum(), nil
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
			url:            "https://music.apple.com/album/abc123",
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
			name:           "異常系: artistのURL",
			url:            "https://open.spotify.com/artist/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "DIFFERENT_SPOTIFY_URL",
		},
		{
			name: "異常系: APIエラー",
			url:  "https://open.spotify.com/album/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Album, error) {
				return nil, errors.New("API error")
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SOMETHING_SPOTIFY_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &mockSpotifyAPIForAlbum{GetAlbumByIDFunc: tt.mockFunc}
			albumUC := usecase.NewAlbumUseCase(mockAPI)
			handler := NewAlbumHandler(albumUC)

			req := httptest.NewRequest(http.MethodGet, "/v1/album/fetch?url="+tt.url, nil)
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
				if result["id"] != "test-album-id" {
					t.Errorf("expected id test-album-id, got %v", result["id"])
				}
				if result["name"] != "Test Album" {
					t.Errorf("expected name Test Album, got %v", result["name"])
				}
				if result["release_date"] != "2024-01-15" {
					t.Errorf("expected release_date 2024-01-15, got %v", result["release_date"])
				}
				// Check tracks
				tracks, ok := result["tracks"].(map[string]interface{})
				if !ok {
					t.Fatal("expected tracks to be object")
				}
				items, ok := tracks["items"].([]interface{})
				if !ok {
					t.Fatal("expected tracks.items to be array")
				}
				if len(items) != 2 {
					t.Errorf("expected 2 tracks, got %d", len(items))
				}
			}
		})
	}
}

func TestAlbumHandler_FetchByURL_EmptyTracks(t *testing.T) {
	mockAPI := &mockSpotifyAPIForAlbum{
		GetAlbumByIDFunc: func(ctx context.Context, id string) (*domain.Album, error) {
			return &domain.Album{
				ID:          "test-id",
				Name:        "Empty Album",
				URL:         "https://open.spotify.com/album/test-id",
				ReleaseDate: "2024-01-01",
				Artists:     []domain.Artist{},
				Tracks:      []domain.SimpleTrack{},
				Genres:      []string{},
				Images:      []domain.Image{},
			}, nil
		},
	}
	albumUC := usecase.NewAlbumUseCase(mockAPI)
	handler := NewAlbumHandler(albumUC)

	req := httptest.NewRequest(http.MethodGet, "/v1/album/fetch?url=https://open.spotify.com/album/abc123", nil)
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

	tracks, ok := result["tracks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected tracks to be object")
	}
	items, ok := tracks["items"].([]interface{})
	if !ok {
		t.Fatal("expected tracks.items to be array")
	}
	if len(items) != 0 {
		t.Errorf("expected 0 tracks, got %d", len(items))
	}
}
