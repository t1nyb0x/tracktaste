package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/testutil"
)

func TestTrackUseCase_FetchByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*testutil.MockSpotifyAPI)
		wantTrack *domain.Track
		wantErr   error
	}{
		{
			name: "正常系: 有効なID",
			id:   "track123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return testutil.CreateTestTrack(id, "Test Track"), nil
				}
			},
			wantTrack: testutil.CreateTestTrack("track123", "Test Track"),
			wantErr:   nil,
		},
		{
			name:      "異常系: 空のID",
			id:        "",
			setupMock: func(m *testutil.MockSpotifyAPI) {},
			wantTrack: nil,
			wantErr:   domain.ErrTrackNotFound,
		},
		{
			name: "異常系: APIエラー",
			id:   "track123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return nil, errors.New("api error")
				}
			},
			wantTrack: nil,
			wantErr:   errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &testutil.MockSpotifyAPI{}
			tt.setupMock(mockAPI)

			uc := NewTrackUseCase(mockAPI)
			got, err := uc.FetchByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("FetchByID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrTrackNotFound {
					if !errors.Is(err, domain.ErrTrackNotFound) {
						t.Errorf("FetchByID() error = %v, want %v", err, tt.wantErr)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FetchByID() unexpected error = %v", err)
				return
			}

			if got.ID != tt.wantTrack.ID || got.Name != tt.wantTrack.Name {
				t.Errorf("FetchByID() = %v, want %v", got, tt.wantTrack)
			}
		})
	}
}

func TestTrackUseCase_Search(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		setupMock func(*testutil.MockSpotifyAPI)
		wantCount int
		wantErr   error
	}{
		{
			name:  "正常系: 有効なクエリ",
			query: "test query",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.SearchTracksFunc = func(ctx context.Context, query string) ([]domain.Track, error) {
					return []domain.Track{
						*testutil.CreateTestTrack("track1", "Track 1"),
						*testutil.CreateTestTrack("track2", "Track 2"),
					}, nil
				}
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:      "異常系: 空クエリ",
			query:     "",
			setupMock: func(m *testutil.MockSpotifyAPI) {},
			wantCount: 0,
			wantErr:   domain.ErrEmptyQuery,
		},
		{
			name:  "正常系: 結果0件",
			query: "no results",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.SearchTracksFunc = func(ctx context.Context, query string) ([]domain.Track, error) {
					return []domain.Track{}, nil
				}
			},
			wantCount: 0,
			wantErr:   nil,
		},
		{
			name:  "異常系: APIエラー",
			query: "test",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.SearchTracksFunc = func(ctx context.Context, query string) ([]domain.Track, error) {
					return nil, errors.New("api error")
				}
			},
			wantCount: 0,
			wantErr:   errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &testutil.MockSpotifyAPI{}
			tt.setupMock(mockAPI)

			uc := NewTrackUseCase(mockAPI)
			got, err := uc.Search(context.Background(), tt.query)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Search() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrEmptyQuery {
					if !errors.Is(err, domain.ErrEmptyQuery) {
						t.Errorf("Search() error = %v, want %v", err, tt.wantErr)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Search() unexpected error = %v", err)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("Search() returned %d tracks, want %d", len(got), tt.wantCount)
			}
		})
	}
}
