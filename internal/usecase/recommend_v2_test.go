package usecase

import (
	"context"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
)

// Mock implementations for testing

type mockDeezerAPI struct {
	tracks map[string]*domain.DeezerTrack
}

func (m *mockDeezerAPI) GetTrackByISRC(ctx context.Context, isrc string) (*domain.DeezerTrack, error) {
	if track, ok := m.tracks[isrc]; ok {
		return track, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockDeezerAPI) SearchTrack(ctx context.Context, title, artist string) (*domain.DeezerTrack, error) {
	return nil, domain.ErrNotFound
}

func (m *mockDeezerAPI) GetTracksByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.DeezerTrack, error) {
	result := make(map[string]*domain.DeezerTrack)
	for _, isrc := range isrcs {
		if track, ok := m.tracks[isrc]; ok {
			result[isrc] = track
		}
	}
	return result, nil
}

type mockMusicBrainzAPI struct {
	recordings map[string]*domain.MBRecording
	artists    map[string]*domain.MBArtist
}

func (m *mockMusicBrainzAPI) GetRecordingByISRC(ctx context.Context, isrc string) (*domain.MBRecording, error) {
	if rec, ok := m.recordings[isrc]; ok {
		return rec, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockMusicBrainzAPI) GetRecordingWithTags(ctx context.Context, mbid string) (*domain.MBRecording, error) {
	return nil, domain.ErrNotFound
}

func (m *mockMusicBrainzAPI) GetArtistWithRelations(ctx context.Context, mbid string) (*domain.MBArtist, error) {
	if artist, ok := m.artists[mbid]; ok {
		return artist, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockMusicBrainzAPI) GetRecordingsByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.MBRecording, error) {
	result := make(map[string]*domain.MBRecording)
	for _, isrc := range isrcs {
		if rec, ok := m.recordings[isrc]; ok {
			result[isrc] = rec
		}
	}
	return result, nil
}

type mockSpotifyAPIV2 struct {
	tracks       map[string]*domain.Track  // by track ID
	tracksByISRC map[string]*domain.Track  // by ISRC
	artists      map[string][]string       // artist ID -> genres
}

func (m *mockSpotifyAPIV2) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	if track, ok := m.tracks[id]; ok {
		return track, nil
	}
	return nil, domain.ErrTrackNotFound
}

func (m *mockSpotifyAPIV2) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	return nil, domain.ErrArtistNotFound
}

func (m *mockSpotifyAPIV2) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	return nil, domain.ErrAlbumNotFound
}

func (m *mockSpotifyAPIV2) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIV2) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	if m.tracksByISRC != nil {
		if track, ok := m.tracksByISRC[isrc]; ok {
			return track, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSpotifyAPIV2) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	return nil, domain.ErrNotFound
}

func (m *mockSpotifyAPIV2) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPIV2) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPIV2) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	if genres, ok := m.artists[artistID]; ok {
		return genres, nil
	}
	return nil, nil
}

func (m *mockSpotifyAPIV2) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, id := range artistIDs {
		if genres, ok := m.artists[id]; ok {
			result[id] = genres
		}
	}
	return result, nil
}

type mockKKBOXAPIV2 struct {
	tracks     map[string]*external.KKBOXTrackInfo
	recommended []external.KKBOXTrackInfo
}

func (m *mockKKBOXAPIV2) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	if track, ok := m.tracks[isrc]; ok {
		return track, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockKKBOXAPIV2) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	return m.recommended, nil
}

func (m *mockKKBOXAPIV2) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	return nil, domain.ErrNotFound
}

func TestRecommendUseCaseV2_GetRecommendations(t *testing.T) {
	isrc := "JPAB12345678"
	trackID := "spotify-track-123"
	artistID := "spotify-artist-123"

	isrc1 := "JPAB00000001"
	isrc2 := "JPAB00000002"

	spotifyAPI := &mockSpotifyAPIV2{
		tracks: map[string]*domain.Track{
			trackID: {
				ID:   trackID,
				Name: "Test Track",
				ISRC: &isrc,
				Artists: []domain.Artist{
					{ID: artistID, Name: "Test Artist"},
				},
			},
		},
		tracksByISRC: map[string]*domain.Track{
			isrc1: {
				ID:   "spotify-rec-1",
				Name: "Recommended 1",
				ISRC: &isrc1,
				Artists: []domain.Artist{
					{ID: "artist-rec-1", Name: "Rec Artist 1"},
				},
				Album: domain.Album{
					ID:   "album-rec-1",
					Name: "Rec Album 1",
					URL:  "https://open.spotify.com/album/album-rec-1",
				},
				URL: "https://open.spotify.com/track/spotify-rec-1",
			},
			isrc2: {
				ID:   "spotify-rec-2",
				Name: "Recommended 2",
				ISRC: &isrc2,
				Artists: []domain.Artist{
					{ID: "artist-rec-2", Name: "Rec Artist 2"},
				},
				Album: domain.Album{
					ID:   "album-rec-2",
					Name: "Rec Album 2",
					URL:  "https://open.spotify.com/album/album-rec-2",
				},
				URL: "https://open.spotify.com/track/spotify-rec-2",
			},
		},
		artists: map[string][]string{
			artistID: {"anime", "jpop"},
		},
	}

	kkboxAPI := &mockKKBOXAPIV2{
		tracks: map[string]*external.KKBOXTrackInfo{
			isrc: {ID: "kkbox-123", Name: "Test Track", ISRC: isrc},
		},
		recommended: []external.KKBOXTrackInfo{
			{ID: "kkbox-rec-1", Name: "Recommended 1", ISRC: "JPAB00000001"},
			{ID: "kkbox-rec-2", Name: "Recommended 2", ISRC: "JPAB00000002"},
		},
	}

	deezerAPI := &mockDeezerAPI{
		tracks: map[string]*domain.DeezerTrack{
			isrc: {
				ID:              123,
				Title:           "Test Track",
				ISRC:            isrc,
				BPM:             175.0,
				DurationSeconds: 245,
				Gain:            -7.2,
			},
			"JPAB00000001": {
				ID:              124,
				Title:           "Recommended 1",
				ISRC:            "JPAB00000001",
				BPM:             180.0,
				DurationSeconds: 240,
				Gain:            -6.8,
			},
			"JPAB00000002": {
				ID:              125,
				Title:           "Recommended 2",
				ISRC:            "JPAB00000002",
				BPM:             120.0,
				DurationSeconds: 300,
				Gain:            -10.0,
			},
		},
	}

	mbAPI := &mockMusicBrainzAPI{
		recordings: map[string]*domain.MBRecording{
			isrc: {
				MBID:       "mb-recording-123",
				Title:      "Test Track",
				ISRC:       isrc,
				Tags:       []domain.MBTag{{Name: "anime", Count: 10}},
				ArtistMBID: "mb-artist-123",
			},
			"JPAB00000001": {
				MBID:       "mb-recording-124",
				Title:      "Recommended 1",
				ISRC:       "JPAB00000001",
				Tags:       []domain.MBTag{{Name: "anime", Count: 8}, {Name: "jpop", Count: 5}},
				ArtistMBID: "mb-artist-124",
			},
		},
		artists: map[string]*domain.MBArtist{
			"mb-artist-123": {
				MBID: "mb-artist-123",
				Name: "Test Artist",
				Tags: []domain.MBTag{{Name: "japanese", Count: 10}},
			},
		},
	}

	uc := NewRecommendUseCaseV2(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	result, err := uc.GetRecommendations(context.Background(), trackID, domain.RecommendModeBalanced, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result should not be nil")
	}

	if result.SeedTrack.ID != trackID {
		t.Errorf("SeedTrack.ID = %s, want %s", result.SeedTrack.ID, trackID)
	}

	if result.Mode != domain.RecommendModeBalanced {
		t.Errorf("Mode = %s, want %s", result.Mode, domain.RecommendModeBalanced)
	}

	if len(result.Items) != 2 {
		t.Errorf("len(Items) = %d, want 2", len(result.Items))
	}

	// Check that results are sorted by score
	if len(result.Items) >= 2 {
		if result.Items[0].FinalScore < result.Items[1].FinalScore {
			t.Error("Items should be sorted by FinalScore descending")
		}
	}
}

func TestRecommendUseCaseV2_GetRecommendations_NoISRC(t *testing.T) {
	trackID := "spotify-track-no-isrc"

	spotifyAPI := &mockSpotifyAPIV2{
		tracks: map[string]*domain.Track{
			trackID: {
				ID:   trackID,
				Name: "Test Track",
				ISRC: nil, // No ISRC
			},
		},
	}

	kkboxAPI := &mockKKBOXAPIV2{}
	deezerAPI := &mockDeezerAPI{}
	mbAPI := &mockMusicBrainzAPI{}

	uc := NewRecommendUseCaseV2(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	result, err := uc.GetRecommendations(context.Background(), trackID, domain.RecommendModeBalanced, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("expected 0 items for track without ISRC, got %d", len(result.Items))
	}
}

func TestRecommendUseCaseV2_MergeTags(t *testing.T) {
	uc := &RecommendUseCaseV2{}

	tests := []struct {
		name          string
		mbTags        []string
		spotifyGenres []string
		expectedLen   int
	}{
		{
			name:          "no duplicates",
			mbTags:        []string{"anime", "jpop"},
			spotifyGenres: []string{"rock", "pop"},
			expectedLen:   4,
		},
		{
			name:          "with duplicates",
			mbTags:        []string{"anime", "jpop"},
			spotifyGenres: []string{"anime", "rock"},
			expectedLen:   3,
		},
		{
			name:          "empty mb tags",
			mbTags:        []string{},
			spotifyGenres: []string{"anime", "jpop"},
			expectedLen:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.mergeTags(tt.mbTags, tt.spotifyGenres)
			if len(result) != tt.expectedLen {
				t.Errorf("mergeTags() returned %d items, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestNewRecommendUseCaseV2(t *testing.T) {
	spotifyAPI := &mockSpotifyAPIV2{}
	kkboxAPI := &mockKKBOXAPIV2{}
	deezerAPI := &mockDeezerAPI{}
	mbAPI := &mockMusicBrainzAPI{}

	uc := NewRecommendUseCaseV2(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	if uc == nil {
		t.Fatal("NewRecommendUseCaseV2 returned nil")
	}
	if uc.calculatorV2 == nil {
		t.Error("calculatorV2 should not be nil")
	}
	if uc.genreMatcher == nil {
		t.Error("genreMatcher should not be nil")
	}
}
