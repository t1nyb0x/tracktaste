package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/testutil"
)

func TestArtistUseCase_FetchByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(*testutil.MockSpotifyAPI)
		wantArtist *domain.Artist
		wantErr    error
	}{
		{
			name: "正常系: 有効なID",
			id:   "artist123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetArtistByIDFunc = func(ctx context.Context, id string) (*domain.Artist, error) {
					return testutil.CreateTestArtist(id, "Test Artist"), nil
				}
			},
			wantArtist: testutil.CreateTestArtist("artist123", "Test Artist"),
			wantErr:    nil,
		},
		{
			name:       "異常系: 空のID",
			id:         "",
			setupMock:  func(m *testutil.MockSpotifyAPI) {},
			wantArtist: nil,
			wantErr:    domain.ErrArtistNotFound,
		},
		{
			name: "異常系: APIエラー",
			id:   "artist123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetArtistByIDFunc = func(ctx context.Context, id string) (*domain.Artist, error) {
					return nil, errors.New("api error")
				}
			},
			wantArtist: nil,
			wantErr:    errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &testutil.MockSpotifyAPI{}
			tt.setupMock(mockAPI)

			uc := NewArtistUseCase(mockAPI)
			got, err := uc.FetchByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("FetchByID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrArtistNotFound {
					if !errors.Is(err, domain.ErrArtistNotFound) {
						t.Errorf("FetchByID() error = %v, want %v", err, tt.wantErr)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FetchByID() unexpected error = %v", err)
				return
			}

			if got.ID != tt.wantArtist.ID || got.Name != tt.wantArtist.Name {
				t.Errorf("FetchByID() = %v, want %v", got, tt.wantArtist)
			}
		})
	}
}
