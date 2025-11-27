package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/testutil"
)

func TestAlbumUseCase_FetchByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*testutil.MockSpotifyAPI)
		wantAlbum *domain.Album
		wantErr   error
	}{
		{
			name: "正常系: 有効なID",
			id:   "album123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetAlbumByIDFunc = func(ctx context.Context, id string) (*domain.Album, error) {
					return testutil.CreateTestAlbum(id, "Test Album"), nil
				}
			},
			wantAlbum: testutil.CreateTestAlbum("album123", "Test Album"),
			wantErr:   nil,
		},
		{
			name:      "異常系: 空のID",
			id:        "",
			setupMock: func(m *testutil.MockSpotifyAPI) {},
			wantAlbum: nil,
			wantErr:   domain.ErrAlbumNotFound,
		},
		{
			name: "異常系: APIエラー",
			id:   "album123",
			setupMock: func(m *testutil.MockSpotifyAPI) {
				m.GetAlbumByIDFunc = func(ctx context.Context, id string) (*domain.Album, error) {
					return nil, errors.New("api error")
				}
			},
			wantAlbum: nil,
			wantErr:   errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &testutil.MockSpotifyAPI{}
			tt.setupMock(mockAPI)

			uc := NewAlbumUseCase(mockAPI)
			got, err := uc.FetchByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("FetchByID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrAlbumNotFound {
					if !errors.Is(err, domain.ErrAlbumNotFound) {
						t.Errorf("FetchByID() error = %v, want %v", err, tt.wantErr)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FetchByID() unexpected error = %v", err)
				return
			}

			if got.ID != tt.wantAlbum.ID || got.Name != tt.wantAlbum.Name {
				t.Errorf("FetchByID() = %v, want %v", got, tt.wantAlbum)
			}
		})
	}
}
