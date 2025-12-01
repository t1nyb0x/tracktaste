package v1

import (
	"context"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/testutil"
)

func TestRecommendUseCase_GetRecommendations(t *testing.T) {
	tests := []struct {
		name      string
		trackID   string
		mode      domain.RecommendMode
		limit     int
		setupMock func(*testutil.MockSpotifyAPI, *testutil.MockKKBOXAPI)
		wantErr   bool
		wantItems int
	}{
		{
			name:    "successful recommendation",
			trackID: "track1",
			mode:    domain.RecommendModeBalanced,
			limit:   10,
			setupMock: func(spotify *testutil.MockSpotifyAPI, kkbox *testutil.MockKKBOXAPI) {
				isrc := "JPTEST12345"
				spotify.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return &domain.Track{
						ID:   id,
						Name: "Test Track",
						ISRC: &isrc,
						Artists: []domain.Artist{
							{ID: "artist1", Name: "Test Artist"},
						},
						Album: domain.Album{ID: "album1", Name: "Test Album"},
					}, nil
				}
				spotify.GetAudioFeaturesFunc = func(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
					return &domain.AudioFeatures{
						TrackID: trackID,
						Tempo:   128.0,
						Energy:  0.8,
						Valence: 0.6,
					}, nil
				}
				spotify.GetArtistGenresFunc = func(ctx context.Context, artistID string) ([]string, error) {
					return []string{"anime"}, nil
				}
				spotify.GetRecommendationsFunc = func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
					return []domain.Track{
						{ID: "rec1", Name: "Recommended 1", Artists: []domain.Artist{{ID: "a1", Name: "Artist 1"}}, Album: domain.Album{ID: "alb1"}},
						{ID: "rec2", Name: "Recommended 2", Artists: []domain.Artist{{ID: "a2", Name: "Artist 2"}}, Album: domain.Album{ID: "alb2"}},
					}, nil
				}
				spotify.GetAudioFeaturesBatchFunc = func(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
					features := make([]domain.AudioFeatures, len(trackIDs))
					for i, id := range trackIDs {
						features[i] = domain.AudioFeatures{
							TrackID: id,
							Tempo:   130.0,
							Energy:  0.75,
							Valence: 0.65,
						}
					}
					return features, nil
				}
				spotify.GetArtistGenresBatchFunc = func(ctx context.Context, artistIDs []string) (map[string][]string, error) {
					result := make(map[string][]string)
					for _, id := range artistIDs {
						result[id] = []string{"anime"}
					}
					return result, nil
				}
				kkbox.SearchByISRCFunc = func(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
					return &external.KKBOXTrackInfo{ID: "kkbox1", Name: "Test", ISRC: isrc}, nil
				}
				kkbox.GetRecommendedTracksFunc = func(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
					return []external.KKBOXTrackInfo{
						{ID: "kkbox2", Name: "KKBOX Rec 1", ISRC: "JPKKBOX00001"},
					}, nil
				}
				spotify.SearchByISRCFunc = func(ctx context.Context, isrc string) (*domain.Track, error) {
					return &domain.Track{
						ID:      "spotify_from_kkbox",
						Name:    "From KKBOX",
						Artists: []domain.Artist{{ID: "a3", Name: "Artist 3"}},
						Album:   domain.Album{ID: "alb3"},
					}, nil
				}
			},
			wantErr:   false,
			wantItems: 3, // 2 from Spotify + 1 from KKBOX
		},
		{
			name:    "no candidates",
			trackID: "track2",
			mode:    domain.RecommendModeBalanced,
			limit:   10,
			setupMock: func(spotify *testutil.MockSpotifyAPI, kkbox *testutil.MockKKBOXAPI) {
				spotify.GetTrackByIDFunc = func(ctx context.Context, id string) (*domain.Track, error) {
					return &domain.Track{
						ID:   id,
						Name: "Test Track",
						Artists: []domain.Artist{
							{ID: "artist1", Name: "Test Artist"},
						},
						Album: domain.Album{ID: "album1", Name: "Test Album"},
					}, nil
				}
				spotify.GetAudioFeaturesFunc = func(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
					return nil, nil
				}
				spotify.GetArtistGenresFunc = func(ctx context.Context, artistID string) ([]string, error) {
					return nil, nil
				}
				spotify.GetRecommendationsFunc = func(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
					return []domain.Track{}, nil
				}
			},
			wantErr:   false,
			wantItems: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotifyMock := &testutil.MockSpotifyAPI{}
			kkboxMock := &testutil.MockKKBOXAPI{}
			tt.setupMock(spotifyMock, kkboxMock)

			uc := NewRecommendUseCase(spotifyMock, kkboxMock)
			result, err := uc.GetRecommendations(context.Background(), tt.trackID, tt.mode, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecommendations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(result.Items) != tt.wantItems {
					t.Errorf("GetRecommendations() items = %d, want %d", len(result.Items), tt.wantItems)
				}
				if result.Mode != tt.mode {
					t.Errorf("GetRecommendations() mode = %v, want %v", result.Mode, tt.mode)
				}
			}
		})
	}
}

func TestRemoveDuplicateTracks(t *testing.T) {
	isrc1 := "ISRC001"
	isrc2 := "ISRC002"

	tests := []struct {
		name   string
		tracks []domain.Track
		seedID string
		want   int
	}{
		{
			name: "remove duplicates by ID",
			tracks: []domain.Track{
				{ID: "t1", Name: "Track 1"},
				{ID: "t1", Name: "Track 1 Dup"},
				{ID: "t2", Name: "Track 2"},
			},
			seedID: "seed",
			want:   2,
		},
		{
			name: "remove duplicates by ISRC",
			tracks: []domain.Track{
				{ID: "t1", Name: "Track 1", ISRC: &isrc1},
				{ID: "t2", Name: "Track 1 Diff ID", ISRC: &isrc1},
				{ID: "t3", Name: "Track 2", ISRC: &isrc2},
			},
			seedID: "seed",
			want:   2,
		},
		{
			name: "remove seed track",
			tracks: []domain.Track{
				{ID: "seed", Name: "Seed Track"},
				{ID: "t1", Name: "Track 1"},
			},
			seedID: "seed",
			want:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicateTracks(tt.tracks, tt.seedID)
			if len(got) != tt.want {
				t.Errorf("removeDuplicateTracks() = %d tracks, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCollectUniqueArtistIDs(t *testing.T) {
	tracks := []domain.Track{
		{ID: "t1", Artists: []domain.Artist{{ID: "a1"}, {ID: "a2"}}},
		{ID: "t2", Artists: []domain.Artist{{ID: "a1"}, {ID: "a3"}}},
		{ID: "t3", Artists: []domain.Artist{{ID: "a2"}}},
	}

	got := collectUniqueArtistIDs(tracks)
	if len(got) != 3 {
		t.Errorf("collectUniqueArtistIDs() = %d, want 3", len(got))
	}

	// Check all unique IDs are present
	idSet := make(map[string]bool)
	for _, id := range got {
		idSet[id] = true
	}
	for _, expected := range []string{"a1", "a2", "a3"} {
		if !idSet[expected] {
			t.Errorf("collectUniqueArtistIDs() missing %s", expected)
		}
	}
}
