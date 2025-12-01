package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	usecasev1 "github.com/t1nyb0x/tracktaste/internal/usecase/v1"
)

// mockSpotifyAPI for handler tests
type mockSpotifyAPI struct {
	GetTrackByIDFunc          func(ctx context.Context, id string) (*domain.Track, error)
	SearchTracksFunc          func(ctx context.Context, query string) ([]domain.Track, error)
	SearchByISRCFunc          func(ctx context.Context, isrc string) (*domain.Track, error)
	GetArtistByIDFunc         func(ctx context.Context, id string) (*domain.Artist, error)
	GetAlbumByIDFunc          func(ctx context.Context, id string) (*domain.Album, error)
	GetAudioFeaturesFunc      func(ctx context.Context, trackID string) (*domain.AudioFeatures, error)
	GetAudioFeaturesBatchFunc func(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error)
	GetRecommendationsFunc    func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error)
	GetArtistGenresFunc       func(ctx context.Context, artistID string) ([]string, error)
	GetArtistGenresBatchFunc  func(ctx context.Context, artistIDs []string) (map[string][]string, error)
}

func (m *mockSpotifyAPI) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	if m.GetTrackByIDFunc != nil {
		return m.GetTrackByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	if m.SearchTracksFunc != nil {
		return m.SearchTracksFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	if m.GetArtistByIDFunc != nil {
		return m.GetArtistByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	if m.GetAlbumByIDFunc != nil {
		return m.GetAlbumByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	if m.SearchByISRCFunc != nil {
		return m.SearchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	if m.GetAudioFeaturesFunc != nil {
		return m.GetAudioFeaturesFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	if m.GetAudioFeaturesBatchFunc != nil {
		return m.GetAudioFeaturesBatchFunc(ctx, trackIDs)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	if m.GetRecommendationsFunc != nil {
		return m.GetRecommendationsFunc(ctx, params)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	if m.GetArtistGenresFunc != nil {
		return m.GetArtistGenresFunc(ctx, artistID)
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	if m.GetArtistGenresBatchFunc != nil {
		return m.GetArtistGenresBatchFunc(ctx, artistIDs)
	}
	return nil, nil
}

var _ external.SpotifyAPI = (*mockSpotifyAPI)(nil)

// mockKKBOXAPI for handler tests
type mockKKBOXAPI struct {
	SearchByISRCFunc         func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error)
	GetRecommendedTracksFunc func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error)
	GetTrackDetailFunc       func(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error)
}

func (m *mockKKBOXAPI) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	if m.SearchByISRCFunc != nil {
		return m.SearchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *mockKKBOXAPI) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	if m.GetRecommendedTracksFunc != nil {
		return m.GetRecommendedTracksFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *mockKKBOXAPI) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	if m.GetTrackDetailFunc != nil {
		return m.GetTrackDetailFunc(ctx, trackID)
	}
	return nil, nil
}

var _ external.KKBOXAPI = (*mockKKBOXAPI)(nil)

func createTestTrack() *domain.Track {
	popularity := 80
	isrc := "JPSO00123456"
	return &domain.Track{
		ID:          "test-track-id",
		Name:        "Test Track",
		URL:         "https://open.spotify.com/track/test-track-id",
		ISRC:        &isrc,
		Popularity:  &popularity,
		TrackNumber: 1,
		DiscNumber:  1,
		Album: domain.Album{
			ID:          "test-album-id",
			Name:        "Test Album",
			URL:         "https://open.spotify.com/album/test-album-id",
			ReleaseDate: "2024-01-01",
			Images: []domain.Image{
				{URL: "https://example.com/image.jpg", Height: 300, Width: 300},
			},
			Artists: []domain.Artist{
				{ID: "artist-1", Name: "Test Artist", URL: "https://open.spotify.com/artist/artist-1"},
			},
		},
		Artists: []domain.Artist{
			{ID: "artist-1", Name: "Test Artist", URL: "https://open.spotify.com/artist/artist-1"},
		},
	}
}

func TestTrackHandler_FetchByURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockFunc       func(ctx context.Context, id string) (*domain.Track, error)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "正常系: 有効なURL",
			url:  "https://open.spotify.com/track/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Track, error) {
				return createTestTrack(), nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "正常系: intl-ja付きURL",
			url:  "https://open.spotify.com/intl-ja/track/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Track, error) {
				return createTestTrack(), nil
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
			url:            "https://music.apple.com/track/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "NOT_SPOTIFY_URL",
		},
		{
			name:           "異常系: artistのURL",
			url:            "https://open.spotify.com/artist/abc123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "DIFFERENT_SPOTIFY_URL",
		},
		{
			name: "異常系: APIエラー",
			url:  "https://open.spotify.com/track/abc123",
			mockFunc: func(ctx context.Context, id string) (*domain.Track, error) {
				return nil, errors.New("API error")
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SOMETHING_SPOTIFY_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &mockSpotifyAPI{GetTrackByIDFunc: tt.mockFunc}
			trackUC := usecasev1.NewTrackUseCase(mockAPI)
			similarUC := usecasev1.NewSimilarTracksUseCase(mockAPI, &mockKKBOXAPI{})
			handler := NewTrackHandler(trackUC, similarUC)

			req := httptest.NewRequest(http.MethodGet, "/v1/track/fetch?url="+tt.url, nil)
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
				if result["id"] != "test-track-id" {
					t.Errorf("expected id test-track-id, got %v", result["id"])
				}
			}
		})
	}
}

func TestTrackHandler_Search(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockFunc       func(ctx context.Context, query string) ([]domain.Track, error)
		expectedStatus int
		expectedCode   string
		expectedCount  int
	}{
		{
			name:  "正常系: 検索結果あり",
			query: "test query",
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				track := createTestTrack()
				return []domain.Track{*track, *track}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:  "正常系: 検索結果0件",
			query: "no results",
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return []domain.Track{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "異常系: 空クエリ",
			query:          "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "EMPTY_QUERY",
		},
		{
			name:  "異常系: APIエラー",
			query: "test",
			mockFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
				return nil, errors.New("API error")
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SOMETHING_SPOTIFY_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &mockSpotifyAPI{SearchTracksFunc: tt.mockFunc}
			trackUC := usecasev1.NewTrackUseCase(mockAPI)
			similarUC := usecasev1.NewSimilarTracksUseCase(mockAPI, &mockKKBOXAPI{})
			handler := NewTrackHandler(trackUC, similarUC)

			req := httptest.NewRequest(http.MethodGet, "/v1/track/search?q="+url.QueryEscape(tt.query), nil)
			rec := httptest.NewRecorder()

			handler.Search(rec, req)

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
				resultObj, ok := resp["result"].(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be object")
				}
				items, ok := resultObj["items"].([]interface{})
				if !ok {
					t.Fatal("expected result.items to be array")
				}
				if len(items) != tt.expectedCount {
					t.Errorf("expected %d results, got %d", tt.expectedCount, len(items))
				}
			}
		})
	}
}

func TestTrackHandler_FetchSimilar(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		setupMock      func() (*mockSpotifyAPI, *mockKKBOXAPI)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "正常系: 類似トラック取得成功",
			url:  "https://open.spotify.com/track/abc123",
			setupMock: func() (*mockSpotifyAPI, *mockKKBOXAPI) {
				isrc := "JPSO00123456"
				spotifyAPI := &mockSpotifyAPI{
					GetTrackByIDFunc: func(ctx context.Context, id string) (*domain.Track, error) {
						return &domain.Track{ID: id, ISRC: &isrc}, nil
					},
					SearchTracksFunc: nil,
				}
				kkboxAPI := &mockKKBOXAPI{
					SearchByISRCFunc: func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
						return &external.KKBOXTrackInfo{ID: "kkbox-track-id", ISRC: isrc}, nil
					},
					GetRecommendedTracksFunc: func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
						return []external.KKBOXTrackInfo{}, nil
					},
				}
				return spotifyAPI, kkboxAPI
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "異常系: 空のURL",
			url:            "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "EMPTY_PARAM",
		},
		{
			name: "異常系: ISRCなし",
			url:  "https://open.spotify.com/track/abc123",
			setupMock: func() (*mockSpotifyAPI, *mockKKBOXAPI) {
				spotifyAPI := &mockSpotifyAPI{
					GetTrackByIDFunc: func(ctx context.Context, id string) (*domain.Track, error) {
						return &domain.Track{ID: id, ISRC: nil}, nil
					},
				}
				return spotifyAPI, &mockKKBOXAPI{}
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "ISRC_NOT_FOUND",
		},
		{
			name: "異常系: KKBOXで見つからない",
			url:  "https://open.spotify.com/track/abc123",
			setupMock: func() (*mockSpotifyAPI, *mockKKBOXAPI) {
				isrc := "JPSO00123456"
				spotifyAPI := &mockSpotifyAPI{
					GetTrackByIDFunc: func(ctx context.Context, id string) (*domain.Track, error) {
						return &domain.Track{ID: id, ISRC: &isrc}, nil
					},
				}
				kkboxAPI := &mockKKBOXAPI{
					SearchByISRCFunc: func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
						return nil, nil
					},
				}
				return spotifyAPI, kkboxAPI
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "KKBOX_TRACK_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var handler *TrackHandler
			if tt.setupMock != nil {
				spotifyAPI, kkboxAPI := tt.setupMock()
				trackUC := usecasev1.NewTrackUseCase(spotifyAPI)
				similarUC := usecasev1.NewSimilarTracksUseCase(spotifyAPI, kkboxAPI)
				handler = NewTrackHandler(trackUC, similarUC)
			} else {
				mockAPI := &mockSpotifyAPI{}
				trackUC := usecasev1.NewTrackUseCase(mockAPI)
				similarUC := usecasev1.NewSimilarTracksUseCase(mockAPI, &mockKKBOXAPI{})
				handler = NewTrackHandler(trackUC, similarUC)
			}

			req := httptest.NewRequest(http.MethodGet, "/v1/track/similar?url="+tt.url, nil)
			rec := httptest.NewRecorder()

			handler.FetchSimilar(rec, req)

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
		})
	}
}
