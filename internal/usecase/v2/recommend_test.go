package v2

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

func (m *mockMusicBrainzAPI) GetArtistRecordings(ctx context.Context, artistMBID string, limit int) ([]domain.MBRecording, error) {
	return []domain.MBRecording{}, nil
}

type mockSpotifyAPI struct {
	tracks       map[string]*domain.Track
	tracksByISRC map[string]*domain.Track
	artists      map[string][]string
}

func (m *mockSpotifyAPI) GetTrackByID(ctx context.Context, id string) (*domain.Track, error) {
	if track, ok := m.tracks[id]; ok {
		return track, nil
	}
	return nil, domain.ErrTrackNotFound
}

func (m *mockSpotifyAPI) GetArtistByID(ctx context.Context, id string) (*domain.Artist, error) {
	return nil, domain.ErrArtistNotFound
}

func (m *mockSpotifyAPI) GetAlbumByID(ctx context.Context, id string) (*domain.Album, error) {
	return nil, domain.ErrAlbumNotFound
}

func (m *mockSpotifyAPI) SearchTracks(ctx context.Context, query string) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPI) SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error) {
	if m.tracksByISRC != nil {
		if track, ok := m.tracksByISRC[isrc]; ok {
			return track, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSpotifyAPI) GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error) {
	return nil, domain.ErrNotFound
}

func (m *mockSpotifyAPI) GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error) {
	return nil, nil
}

func (m *mockSpotifyAPI) GetRecommendations(ctx context.Context, params external.RecommendationParams) ([]domain.Track, error) {
	return nil, nil
}

func (m *mockSpotifyAPI) GetArtistGenres(ctx context.Context, artistID string) ([]string, error) {
	if genres, ok := m.artists[artistID]; ok {
		return genres, nil
	}
	return nil, nil
}

func (m *mockSpotifyAPI) GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, id := range artistIDs {
		if genres, ok := m.artists[id]; ok {
			result[id] = genres
		}
	}
	return result, nil
}

type mockKKBOXAPI struct {
	tracks          map[string]*external.KKBOXTrackInfo
	recommended     []external.KKBOXTrackInfo
	returnNilOnMiss bool // if true, return (nil, nil) instead of (nil, error) when track not found
}

func (m *mockKKBOXAPI) SearchByISRC(ctx context.Context, isrc string) (*external.KKBOXTrackInfo, error) {
	if track, ok := m.tracks[isrc]; ok {
		return track, nil
	}
	if m.returnNilOnMiss {
		return nil, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockKKBOXAPI) GetRecommendedTracks(ctx context.Context, trackID string) ([]external.KKBOXTrackInfo, error) {
	return m.recommended, nil
}

func (m *mockKKBOXAPI) GetTrackDetail(ctx context.Context, trackID string) (*external.KKBOXTrackInfo, error) {
	return nil, domain.ErrNotFound
}

func TestRecommendUseCase_GetRecommendations(t *testing.T) {
	isrc := "JPAB12345678"
	trackID := "spotify-track-123"
	artistID := "spotify-artist-123"

	isrc1 := "JPAB00000001"
	isrc2 := "JPAB00000002"

	spotifyAPI := &mockSpotifyAPI{
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

	kkboxAPI := &mockKKBOXAPI{
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

	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

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

func TestRecommendUseCase_GetRecommendations_NoISRC(t *testing.T) {
	trackID := "spotify-track-no-isrc"

	spotifyAPI := &mockSpotifyAPI{
		tracks: map[string]*domain.Track{
			trackID: {
				ID:   trackID,
				Name: "Test Track",
				ISRC: nil,
			},
		},
	}

	kkboxAPI := &mockKKBOXAPI{}
	deezerAPI := &mockDeezerAPI{}
	mbAPI := &mockMusicBrainzAPI{}

	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	result, err := uc.GetRecommendations(context.Background(), trackID, domain.RecommendModeBalanced, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("expected 0 items for track without ISRC, got %d", len(result.Items))
	}
}

func TestRecommendUseCase_MergeTags(t *testing.T) {
	uc := &RecommendUseCase{}

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

func TestNewRecommendUseCase(t *testing.T) {
	spotifyAPI := &mockSpotifyAPI{}
	kkboxAPI := &mockKKBOXAPI{}
	deezerAPI := &mockDeezerAPI{}
	mbAPI := &mockMusicBrainzAPI{}

	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	if uc == nil {
		t.Fatal("NewRecommendUseCase returned nil")
	}
	if uc.calculator == nil {
		t.Error("calculator should not be nil")
	}
	if uc.genreMatcher == nil {
		t.Error("genreMatcher should not be nil")
	}
}

func TestRecommendUseCase_GetRecommendations_KKBOXTrackNotFound(t *testing.T) {
	// Test case for the specific issue: when KKBOX returns nil for SearchByISRC
	// This should not panic, but should return empty recommendations
	isrc := "USRC17607839"
	trackID := "087sGVlyEXq6bDpgnGx78E"
	artistID := "spotify-artist-123"

	spotifyAPI := &mockSpotifyAPI{
		tracks: map[string]*domain.Track{
			trackID: {
				ID:   trackID,
				Name: "A SEAZER 絶対運命黙示録・完全版",
				ISRC: &isrc,
				Artists: []domain.Artist{
					{ID: artistID, Name: "J.A. Seazer"},
				},
			},
		},
		artists: map[string][]string{
			artistID: {"anime", "jpop"},
		},
	}

	// KKBOX returns nil for this ISRC (not found in KKBOX catalog)
	kkboxAPI := &mockKKBOXAPI{
		tracks:          map[string]*external.KKBOXTrackInfo{},
		returnNilOnMiss: true, // Simulate real KKBOX behavior: returns (nil, nil) when not found
	}

	deezerAPI := &mockDeezerAPI{
		tracks: map[string]*domain.DeezerTrack{
			isrc: {
				ID:              123,
				Title:           "A SEAZER 絶対運命黙示録・完全版",
				ISRC:            isrc,
				BPM:             120.0,
				DurationSeconds: 300,
				Gain:            -7.0,
			},
		},
	}

	mbAPI := &mockMusicBrainzAPI{
		recordings: map[string]*domain.MBRecording{
			isrc: {
				MBID:       "mb-recording-123",
				Title:      "A SEAZER 絶対運命黙示録・完全版",
				ISRC:       isrc,
				Tags:       []domain.MBTag{{Name: "anime", Count: 10}},
				ArtistMBID: "mb-artist-123",
			},
		},
		artists: map[string]*domain.MBArtist{
			"mb-artist-123": {
				MBID: "mb-artist-123",
				Name: "J.A. Seazer",
				Tags: []domain.MBTag{{Name: "japanese", Count: 10}},
			},
		},
	}

	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, mbAPI)

	// This should not panic even though KKBOX returns nil
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

	// Should return empty recommendations since KKBOX has no data
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items when KKBOX returns nil, got %d", len(result.Items))
	}
}

func TestSanitizeSearchQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes Japanese brackets",
			input: "聖戦と死神 第1部「銀色の死神」",
			want:  "聖戦と死神 第1部 銀色の死神",
		},
		{
			name:  "removes parentheses",
			input: "Track (Remix)",
			want:  "Track Remix",
		},
		{
			name:  "removes special symbols",
			input: "Track～Version",
			want:  "Track Version",
		},
		{
			name:  "collapses multiple spaces",
			input: "Track    Name",
			want:  "Track Name",
		},
		{
			name:  "preserves normal text",
			input: "Normal Track Name",
			want:  "Normal Track Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeSearchQuery(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeSearchQuery(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSimplifyTrackName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes parenthesized suffix",
			input: "Track Name (feat. Artist)",
			want:  "Track Name",
		},
		{
			name:  "removes Japanese bracketed suffix",
			input: "Chronicle 2nd 聖戦と死神 第1部「銀色の死神」 ～戦場を駈ける者～",
			want:  "Chronicle 2nd 聖戦と死神 第1部 銀色の死神 戦場を駈ける者",
		},
		{
			name:  "removes remix indicator",
			input: "Track Name - Remix Version",
			want:  "Track Name",
		},
		{
			name:  "preserves short names",
			input: "AB",
			want:  "AB",
		},
		{
			name:  "handles normal names",
			input: "Normal Track",
			want:  "Normal Track",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyTrackName(tt.input)
			if got != tt.want {
				t.Errorf("simplifyTrackName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFuzzyMatchArtist(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "exact match",
			a:    "Artist Name",
			b:    "Artist Name",
			want: true,
		},
		{
			name: "case insensitive",
			a:    "ARTIST NAME",
			b:    "artist name",
			want: true,
		},
		{
			name: "one contains other",
			a:    "The Artist",
			b:    "Artist",
			want: true,
		},
		{
			name: "ampersand normalization",
			a:    "Artist & Band",
			b:    "Artist and Band",
			want: true,
		},
		{
			name: "different artists",
			a:    "Artist A",
			b:    "Artist B",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyMatchArtist(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("fuzzyMatchArtist(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
