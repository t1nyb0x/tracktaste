package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/usecase"
)

// mockSpotifyAPI for recommend handler tests
type mockRecommendSpotifyAPI struct {
	getTrackByIDFunc          func(ctx context.Context, id string) (*domain.Track, error)
	getArtistByIDFunc         func(ctx context.Context, id string) (*domain.Artist, error)
	getAlbumByIDFunc          func(ctx context.Context, id string) (*domain.Album, error)
	searchTracksFunc          func(ctx context.Context, query string) ([]domain.Track, error)
	searchByISRCFunc          func(ctx context.Context, isrc string) (*domain.Track, error)
	getAudioFeaturesFunc      func(ctx context.Context, trackID string) (*domain.AudioFeatures, error)
	getAudioFeaturesBatchFunc func(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error)
	getRecommendationsFunc    func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error)
	getArtistGenresFunc       func(ctx context.Context, artistID string) ([]string, error)
	getArtistGenresBatchFunc  func(ctx context.Context, artistIDs []string) (map[string][]string, error)
}

func (m *mockRecommendSpotifyAPI) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	if m.getTrackByIDFunc != nil {
		return m.getTrackByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	if m.getArtistByIDFunc != nil {
		return m.getArtistByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	if m.getAlbumByIDFunc != nil {
		return m.getAlbumByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	if m.searchTracksFunc != nil {
		return m.searchTracksFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	if m.searchByISRCFunc != nil {
		return m.searchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	if m.getAudioFeaturesFunc != nil {
		return m.getAudioFeaturesFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	if m.getAudioFeaturesBatchFunc != nil {
		return m.getAudioFeaturesBatchFunc(ctx, trackIDs)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	if m.getRecommendationsFunc != nil {
		return m.getRecommendationsFunc(ctx, params)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	if m.getArtistGenresFunc != nil {
		return m.getArtistGenresFunc(ctx, artistID)
	}
	return nil, nil
}

func (m *mockRecommendSpotifyAPI) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	if m.getArtistGenresBatchFunc != nil {
		return m.getArtistGenresBatchFunc(ctx, artistIDs)
	}
	return nil, nil
}

// mockKKBOXAPI for recommend handler tests
type mockRecommendKKBOXAPI struct {
	searchByISRCFunc         func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error)
	getRecommendedTracksFunc func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error)
	getTrackDetailFunc       func(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error)
}

func (m *mockRecommendKKBOXAPI) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	if m.searchByISRCFunc != nil {
		return m.searchByISRCFunc(ctx, isrc)
	}
	return nil, nil
}

func (m *mockRecommendKKBOXAPI) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	if m.getRecommendedTracksFunc != nil {
		return m.getRecommendedTracksFunc(ctx, trackID)
	}
	return nil, nil
}

func (m *mockRecommendKKBOXAPI) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	if m.getTrackDetailFunc != nil {
		return m.getTrackDetailFunc(ctx, trackID)
	}
	return nil, nil
}

func TestRecommendHandler_FetchRecommendations(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		setupMock      func(*mockRecommendSpotifyAPI, *mockRecommendKKBOXAPI)
		wantStatusCode int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful recommendation",
			url:  "/v1/track/recommend?url=https://open.spotify.com/track/abc123&mode=balanced&limit=10",
			setupMock: func(spotify *mockRecommendSpotifyAPI, kkbox *mockRecommendKKBOXAPI) {
				spotify.getTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return &domain.Track{
						ID:      id,
						Name:    "Test Track",
						Artists: []domain.Artist{{ID: "a1", Name: "Artist 1"}},
						Album:   domain.Album{ID: "alb1", Name: "Album 1"},
					}, nil
				}
				spotify.getAudioFeaturesFunc = func(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
					return &domain.AudioFeatures{
						TrackID: trackID,
						Tempo:   128.0,
						Energy:  0.8,
					}, nil
				}
				spotify.getArtistGenresFunc = func(ctx context.Context, artistID string) ([]string, error) {
					return []string{"j-pop"}, nil
				}
				spotify.getRecommendationsFunc = func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
					return []domain.Track{
						{ID: "rec1", Name: "Rec 1", Artists: []domain.Artist{{ID: "a2"}}, Album: domain.Album{ID: "alb2"}},
					}, nil
				}
				spotify.getAudioFeaturesBatchFunc = func(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
					return []domain.AudioFeatures{{TrackID: "rec1", Tempo: 130, Energy: 0.75}}, nil
				}
				spotify.getArtistGenresBatchFunc = func(ctx context.Context, artistIDs []string) (map[string][]string, error) {
					return map[string][]string{"a2": {"j-pop"}}, nil
				}
			},
			wantStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Status int `json:"status"`
					Result struct {
						SeedTrack struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"seed_track"`
						Items []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"items"`
						Mode string `json:"mode"`
					} `json:"result"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Result.SeedTrack.ID != "abc123" {
					t.Errorf("Seed track ID = %v, want abc123", resp.Result.SeedTrack.ID)
				}
				if resp.Result.Mode != "balanced" {
					t.Errorf("Mode = %v, want balanced", resp.Result.Mode)
				}
			},
		},
		{
			name:           "missing url parameter",
			url:            "/v1/track/recommend",
			setupMock:      func(spotify *mockRecommendSpotifyAPI, kkbox *mockRecommendKKBOXAPI) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "invalid url format",
			url:            "/v1/track/recommend?url=not-a-spotify-url",
			setupMock:      func(spotify *mockRecommendSpotifyAPI, kkbox *mockRecommendKKBOXAPI) {},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotifyMock := &mockRecommendSpotifyAPI{}
			kkboxMock := &mockRecommendKKBOXAPI{}
			tt.setupMock(spotifyMock, kkboxMock)

			uc := usecase.NewRecommendUseCase(spotifyMock, kkboxMock)
			h := NewRecommendHandler(uc)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			h.FetchRecommendations(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
