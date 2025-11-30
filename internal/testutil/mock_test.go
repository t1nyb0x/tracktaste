package testutil

import (
"context"
"errors"
"testing"

"github.com/t1nyb0x/tracktaste/internal/domain"
"github.com/t1nyb0x/tracktaste/internal/port/external"
)

func TestMockSpotifyAPI_GetTrackByID(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *MockSpotifyAPI
		expectedID  string
		expectedErr bool
	}{
		{
			name: "関数が設定されている場合",
			setupMock: func() *MockSpotifyAPI {
				return &MockSpotifyAPI{
					GetTrackByIDFunc: func(ctx context.Context, id string) (*domain.Track, error) {
						return &domain.Track{ID: id, Name: "Test Track"}, nil
					},
				}
			},
			expectedID:  "track123",
			expectedErr: false,
		},
		{
			name: "関数が設定されていない場合",
			setupMock: func() *MockSpotifyAPI {
				return &MockSpotifyAPI{}
			},
			expectedID:  "",
			expectedErr: false,
		},
		{
			name: "エラーを返す場合",
			setupMock: func() *MockSpotifyAPI {
				return &MockSpotifyAPI{
					GetTrackByIDFunc: func(ctx context.Context, id string) (*domain.Track, error) {
						return nil, errors.New("API error")
					},
				}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
mock := tt.setupMock()
			track, err := mock.GetTrackByID(context.Background(), "track123")

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.expectedID != "" && (track == nil || track.ID != tt.expectedID) {
				t.Errorf("expected ID '%s', got '%v'", tt.expectedID, track)
			}
		})
	}
}

func TestMockSpotifyAPI_GetArtistByID(t *testing.T) {
	mock := &MockSpotifyAPI{
		GetArtistByIDFunc: func(ctx context.Context, id string) (*domain.Artist, error) {
			return &domain.Artist{ID: id, Name: "Test Artist"}, nil
		},
	}

	artist, err := mock.GetArtistByID(context.Background(), "artist123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if artist.ID != "artist123" {
		t.Errorf("expected ID 'artist123', got '%s'", artist.ID)
	}

	emptyMock := &MockSpotifyAPI{}
	artist, err = emptyMock.GetArtistByID(context.Background(), "artist123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if artist != nil {
		t.Error("expected nil artist")
	}
}

func TestMockSpotifyAPI_GetAlbumByID(t *testing.T) {
	mock := &MockSpotifyAPI{
		GetAlbumByIDFunc: func(ctx context.Context, id string) (*domain.Album, error) {
			return &domain.Album{ID: id, Name: "Test Album"}, nil
		},
	}

	album, err := mock.GetAlbumByID(context.Background(), "album123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if album.ID != "album123" {
		t.Errorf("expected ID 'album123', got '%s'", album.ID)
	}

	emptyMock := &MockSpotifyAPI{}
	album, err = emptyMock.GetAlbumByID(context.Background(), "album123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if album != nil {
		t.Error("expected nil album")
	}
}

func TestMockSpotifyAPI_SearchTracks(t *testing.T) {
	mock := &MockSpotifyAPI{
		SearchTracksFunc: func(ctx context.Context, query string) ([]domain.Track, error) {
			return []domain.Track{{ID: "track1", Name: query}}, nil
		},
	}

	tracks, err := mock.SearchTracks(context.Background(), "test query")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tracks) != 1 {
		t.Errorf("expected 1 track, got %d", len(tracks))
	}

	emptyMock := &MockSpotifyAPI{}
	tracks, err = emptyMock.SearchTracks(context.Background(), "test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tracks != nil {
		t.Error("expected nil tracks")
	}
}

func TestMockSpotifyAPI_SearchByISRC(t *testing.T) {
	mock := &MockSpotifyAPI{
		SearchByISRCFunc: func(ctx context.Context, isrc string) (*domain.Track, error) {
			return &domain.Track{ID: "track1", Name: "Found by ISRC"}, nil
		},
	}

	track, err := mock.SearchByISRC(context.Background(), "USRC12345678")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if track == nil || track.ID != "track1" {
		t.Error("expected track to be found")
	}

	emptyMock := &MockSpotifyAPI{}
	track, err = emptyMock.SearchByISRC(context.Background(), "ISRC")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if track != nil {
		t.Error("expected nil track")
	}
}

func TestMockKKBOXAPI_SearchByISRC(t *testing.T) {
	mock := &MockKKBOXAPI{
		SearchByISRCFunc: func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
			return &external.KKBOXTrackInfo{ID: "kkbox123", Name: "KKBOX Track", ISRC: isrc}, nil
		},
	}

	info, err := mock.SearchByISRC(context.Background(), "JPSO00123456")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info.ID != "kkbox123" {
		t.Errorf("expected ID 'kkbox123', got '%s'", info.ID)
	}

	emptyMock := &MockKKBOXAPI{}
	info, err = emptyMock.SearchByISRC(context.Background(), "ISRC")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info != nil {
		t.Error("expected nil info")
	}
}

func TestMockKKBOXAPI_GetRecommendedTracks(t *testing.T) {
	mock := &MockKKBOXAPI{
		GetRecommendedTracksFunc: func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
			return []external.KKBOXTrackInfo{
				{ID: "rec1", Name: "Recommended 1"},
				{ID: "rec2", Name: "Recommended 2"},
			}, nil
		},
	}

	tracks, err := mock.GetRecommendedTracks(context.Background(), "track123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 tracks, got %d", len(tracks))
	}

	emptyMock := &MockKKBOXAPI{}
	tracks, err = emptyMock.GetRecommendedTracks(context.Background(), "track123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tracks != nil {
		t.Error("expected nil tracks")
	}
}

func TestMockKKBOXAPI_GetTrackDetail(t *testing.T) {
	mock := &MockKKBOXAPI{
		GetTrackDetailFunc: func(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
			return &external.KKBOXTrackInfo{ID: trackID, Name: "Detail Track"}, nil
		},
	}

	info, err := mock.GetTrackDetail(context.Background(), "track123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info.ID != "track123" {
		t.Errorf("expected ID 'track123', got '%s'", info.ID)
	}

	emptyMock := &MockKKBOXAPI{}
	info, err = emptyMock.GetTrackDetail(context.Background(), "track123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info != nil {
		t.Error("expected nil info")
	}
}

func TestMockTokenRepository(t *testing.T) {
	t.Run("SaveToken", func(t *testing.T) {
called := false
mock := &MockTokenRepository{
			SaveTokenFunc: func(ctx context.Context, key string, token string, ttlSeconds int) error {
				called = true
				return nil
			},
		}

		err := mock.SaveToken(context.Background(), "key", "token", 3600)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected SaveTokenFunc to be called")
		}

		emptyMock := &MockTokenRepository{}
		err = emptyMock.SaveToken(context.Background(), "key", "token", 3600)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("GetToken", func(t *testing.T) {
mock := &MockTokenRepository{
			GetTokenFunc: func(ctx context.Context, key string) (string, error) {
				return "stored_token", nil
			},
		}

		token, err := mock.GetToken(context.Background(), "key")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if token != "stored_token" {
			t.Errorf("expected 'stored_token', got '%s'", token)
		}

		emptyMock := &MockTokenRepository{}
		token, err = emptyMock.GetToken(context.Background(), "key")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if token != "" {
			t.Errorf("expected empty token, got '%s'", token)
		}
	})

	t.Run("IsTokenValid", func(t *testing.T) {
mock := &MockTokenRepository{
			IsTokenValidFunc: func(ctx context.Context, key string) bool {
				return true
			},
		}

		valid := mock.IsTokenValid(context.Background(), "key")
		if !valid {
			t.Error("expected token to be valid")
		}

		emptyMock := &MockTokenRepository{}
		valid = emptyMock.IsTokenValid(context.Background(), "key")
		if valid {
			t.Error("expected token to be invalid")
		}
	})
}

func TestStringPtr(t *testing.T) {
	s := "test"
	ptr := StringPtr(s)
	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *ptr != s {
		t.Errorf("expected '%s', got '%s'", s, *ptr)
	}
}

func TestIntPtr(t *testing.T) {
	i := 42
	ptr := IntPtr(i)
	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *ptr != i {
		t.Errorf("expected %d, got %d", i, *ptr)
	}
}

func TestCreateTestTrack(t *testing.T) {
	track := CreateTestTrack("track123", "Test Track")

	if track.ID != "track123" {
		t.Errorf("expected ID 'track123', got '%s'", track.ID)
	}
	if track.Name != "Test Track" {
		t.Errorf("expected Name 'Test Track', got '%s'", track.Name)
	}
	if track.URL != "https://open.spotify.com/track/track123" {
		t.Errorf("unexpected URL: %s", track.URL)
	}
}

func TestCreateTestArtist(t *testing.T) {
	artist := CreateTestArtist("artist123", "Test Artist")

	if artist.ID != "artist123" {
		t.Errorf("expected ID 'artist123', got '%s'", artist.ID)
	}
	if artist.Name != "Test Artist" {
		t.Errorf("expected Name 'Test Artist', got '%s'", artist.Name)
	}
}

func TestCreateTestAlbum(t *testing.T) {
	album := CreateTestAlbum("album123", "Test Album")

	if album.ID != "album123" {
		t.Errorf("expected ID 'album123', got '%s'", album.ID)
	}
	if album.Name != "Test Album" {
		t.Errorf("expected Name 'Test Album', got '%s'", album.Name)
	}
}
