package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/testutil"
)

func TestSimilarTracksUseCase_FetchSimilar(t *testing.T) {
	tests := []struct {
		name          string
		trackID       string
		setupSpotify  func(*testutil.MockSpotifyAPI)
		setupKKBOX    func(*testutil.MockKKBOXAPI)
		wantErr       error
		wantItemCount int
	}{
		{
			name:    "正常系: 完全なフロー",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				isrc := "JPTEST12345"
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = &isrc
					return track, nil
				}
				m.SearchByISRCFunc = func(ctx context.Context, isrc string) (*domain.Track, error) {
					pop := 80
					track := testutil.CreateTestTrack("similar1", "Similar Track")
					track.ISRC = testutil.StringPtr(isrc)
					track.Popularity = &pop
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {
				m.SearchByISRCFunc = func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
					return &external.KKBOXTrackInfo{
						ID:   "kkbox123",
						Name: "Test Track",
						ISRC: isrc,
					}, nil
				}
				m.GetRecommendedTracksFunc = func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
					return []external.KKBOXTrackInfo{
						{ID: "rec1", Name: "Rec 1", ISRC: "ISRC001"},
						{ID: "rec2", Name: "Rec 2", ISRC: "ISRC002"},
					}, nil
				}
			},
			wantErr:       nil,
			wantItemCount: 2,
		},
		{
			name:    "正常系: レコメンド0件",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				isrc := "JPTEST12345"
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = &isrc
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {
				m.SearchByISRCFunc = func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
					return &external.KKBOXTrackInfo{
						ID:   "kkbox123",
						Name: "Test Track",
						ISRC: isrc,
					}, nil
				}
				m.GetRecommendedTracksFunc = func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
					return []external.KKBOXTrackInfo{}, nil
				}
			},
			wantErr:       nil,
			wantItemCount: 0,
		},
		{
			name:    "異常系: ISRCなし",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = nil
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {},
			wantErr:    domain.ErrISRCNotFound,
		},
		{
			name:    "異常系: ISRC空文字",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				emptyISRC := ""
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = &emptyISRC
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {},
			wantErr:    domain.ErrISRCNotFound,
		},
		{
			name:    "異常系: KKBOX検索で見つからない",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				isrc := "JPTEST12345"
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = &isrc
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {
				m.SearchByISRCFunc = func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
					return nil, nil
				}
			},
			wantErr: domain.ErrTrackNotFound,
		},
		{
			name:    "異常系: Spotify APIエラー",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return nil, errors.New("spotify api error")
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {},
			wantErr:    errors.New("spotify api error"),
		},
		{
			name:    "異常系: KKBOX検索APIエラー",
			trackID: "track123",
			setupSpotify: func(m *testutil.MockSpotifyAPI) {
				isrc := "JPTEST12345"
				m.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					track := testutil.CreateTestTrack(id, "Test Track")
					track.ISRC = &isrc
					return track, nil
				}
			},
			setupKKBOX: func(m *testutil.MockKKBOXAPI) {
				m.SearchByISRCFunc = func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
					return nil, errors.New("kkbox api error")
				}
			},
			wantErr: errors.New("kkbox api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpotify := &testutil.MockSpotifyAPI{}
			mockKKBOX := &testutil.MockKKBOXAPI{}
			tt.setupSpotify(mockSpotify)
			tt.setupKKBOX(mockKKBOX)

			uc := NewSimilarTracksUseCase(mockSpotify, mockKKBOX)
			got, err := uc.FetchSimilar(context.Background(), tt.trackID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("FetchSimilar() error = nil, wantErr %v", tt.wantErr)
					return
				}
				// Check for domain errors
				if errors.Is(tt.wantErr, domain.ErrISRCNotFound) {
					if !errors.Is(err, domain.ErrISRCNotFound) {
						t.Errorf("FetchSimilar() error = %v, want %v", err, tt.wantErr)
					}
				} else if errors.Is(tt.wantErr, domain.ErrTrackNotFound) {
					if !errors.Is(err, domain.ErrTrackNotFound) {
						t.Errorf("FetchSimilar() error = %v, want %v", err, tt.wantErr)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FetchSimilar() unexpected error = %v", err)
				return
			}

			if len(got.Items) != tt.wantItemCount {
				t.Errorf("FetchSimilar() returned %d items, want %d", len(got.Items), tt.wantItemCount)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		tracks []domain.SimilarTrack
		want   int
	}{
		{
			name: "重複なし",
			tracks: []domain.SimilarTrack{
				{ID: "1", ISRC: testutil.StringPtr("ISRC1")},
				{ID: "2", ISRC: testutil.StringPtr("ISRC2")},
			},
			want: 2,
		},
		{
			name: "ISRC重複あり",
			tracks: []domain.SimilarTrack{
				{ID: "1", ISRC: testutil.StringPtr("ISRC1")},
				{ID: "2", ISRC: testutil.StringPtr("ISRC1")}, // 重複
				{ID: "3", ISRC: testutil.StringPtr("ISRC2")},
			},
			want: 2,
		},
		{
			name: "ID重複（ISRCなし）",
			tracks: []domain.SimilarTrack{
				{ID: "1", ISRC: nil},
				{ID: "1", ISRC: nil}, // 重複
				{ID: "2", ISRC: nil},
			},
			want: 2,
		},
		{
			name:   "空配列",
			tracks: []domain.SimilarTrack{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.tracks)
			if len(got) != tt.want {
				t.Errorf("removeDuplicates() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}
