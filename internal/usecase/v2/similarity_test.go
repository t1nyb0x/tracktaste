package v2

import (
	"math"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

func TestSimilarityCalculator_Calculate(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name      string
		seed      *domain.TrackFeatures
		candidate *domain.TrackFeatures
		wantMin   float64
		wantMax   float64
	}{
		{
			name:      "nil seed",
			seed:      nil,
			candidate: &domain.TrackFeatures{BPM: 120},
			wantMin:   0.5,
			wantMax:   0.5,
		},
		{
			name:      "nil candidate",
			seed:      &domain.TrackFeatures{BPM: 120},
			candidate: nil,
			wantMin:   0.5,
			wantMax:   0.5,
		},
		{
			name: "identical features",
			seed: &domain.TrackFeatures{
				BPM:             175.0,
				DurationSeconds: 245,
				Gain:            -7.2,
				Tags:            []string{"anime", "jpop"},
			},
			candidate: &domain.TrackFeatures{
				BPM:             175.0,
				DurationSeconds: 245,
				Gain:            -7.2,
				Tags:            []string{"anime", "jpop"},
			},
			wantMin: 0.99,
			wantMax: 1.0,
		},
		{
			name: "similar BPM",
			seed: &domain.TrackFeatures{
				BPM:             175.0,
				DurationSeconds: 245,
			},
			candidate: &domain.TrackFeatures{
				BPM:             180.0,
				DurationSeconds: 240,
			},
			wantMin: 0.9,
			wantMax: 1.0,
		},
		{
			name: "very different features",
			seed: &domain.TrackFeatures{
				BPM:             60.0,
				DurationSeconds: 120,
				Gain:            -3.0,
				Tags:            []string{"classical"},
			},
			candidate: &domain.TrackFeatures{
				BPM:             180.0,
				DurationSeconds: 300,
				Gain:            -15.0,
				Tags:            []string{"rock", "metal"},
			},
			wantMin: 0.0,
			wantMax: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Calculate(tt.seed, tt.candidate)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("Calculate() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_bpmSimilarity(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		bpmA    float64
		bpmB    float64
		wantMin float64
		wantMax float64
	}{
		{"identical", 120.0, 120.0, 1.0, 1.0},
		{"close", 120.0, 125.0, 0.95, 1.0},
		{"far apart", 60.0, 180.0, 0.4, 0.6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.bpmSimilarity(tt.bpmA, tt.bpmB)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("bpmSimilarity() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_durationSimilarity(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		durA    int
		durB    int
		wantMin float64
		wantMax float64
	}{
		{"identical", 240, 240, 1.0, 1.0},
		{"close", 240, 250, 0.95, 1.0},
		{"very different", 120, 600, 0.15, 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.durationSimilarity(tt.durA, tt.durB)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("durationSimilarity() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_gainSimilarity(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		gainA   float64
		gainB   float64
		wantMin float64
		wantMax float64
	}{
		{"identical", -7.0, -7.0, 1.0, 1.0},
		{"close", -7.0, -8.0, 0.9, 1.0},
		{"far apart", -3.0, -18.0, 0.15, 0.35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.gainSimilarity(tt.gainA, tt.gainB)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("gainSimilarity() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_tagSimilarity(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		tagsA   []string
		tagsB   []string
		wantMin float64
		wantMax float64
	}{
		{
			name:    "identical",
			tagsA:   []string{"anime", "jpop"},
			tagsB:   []string{"anime", "jpop"},
			wantMin: 1.0,
			wantMax: 1.0,
		},
		{
			name:    "partial overlap",
			tagsA:   []string{"anime", "jpop", "female vocalist"},
			tagsB:   []string{"anime", "rock"},
			wantMin: 0.2,
			wantMax: 0.3,
		},
		{
			name:    "no overlap",
			tagsA:   []string{"classical"},
			tagsB:   []string{"rock", "metal"},
			wantMin: 0.0,
			wantMax: 0.0,
		},
		{
			name:    "both empty",
			tagsA:   []string{},
			tagsB:   []string{},
			wantMin: 0.5,
			wantMax: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.tagSimilarity(tt.tagsA, tt.tagsB)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("tagSimilarity() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_calculateArtistBonus(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name     string
		artistA  *domain.ArtistInfo
		artistB  *domain.ArtistInfo
		expected float64
	}{
		{
			name:     "nil artists",
			artistA:  nil,
			artistB:  nil,
			expected: 1.0,
		},
		{
			name: "same spotify ID",
			artistA: &domain.ArtistInfo{
				SpotifyID: "spotify-123",
			},
			artistB: &domain.ArtistInfo{
				SpotifyID: "spotify-123",
			},
			expected: 1.5,
		},
		{
			name: "same MBID",
			artistA: &domain.ArtistInfo{
				MBID: "mbid-123",
			},
			artistB: &domain.ArtistInfo{
				MBID: "mbid-123",
			},
			expected: 1.5,
		},
		{
			name: "no relation",
			artistA: &domain.ArtistInfo{
				SpotifyID: "spotify-123",
				MBID:      "mbid-123",
			},
			artistB: &domain.ArtistInfo{
				SpotifyID: "spotify-456",
				MBID:      "mbid-456",
			},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.calculateArtistBonus(tt.artistA, tt.artistB)
			if math.Abs(got-tt.expected) > 0.001 {
				t.Errorf("calculateArtistBonus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSimilarityCalculator_MatchReasons(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name      string
		seed      *domain.TrackFeatures
		candidate *domain.TrackFeatures
		wantLen   int
	}{
		{
			name:      "nil features",
			seed:      nil,
			candidate: nil,
			wantLen:   0,
		},
		{
			name: "similar BPM and tags",
			seed: &domain.TrackFeatures{
				BPM:             175.0,
				DurationSeconds: 245,
				Tags:            []string{"anime", "jpop"},
			},
			candidate: &domain.TrackFeatures{
				BPM:             180.0,
				DurationSeconds: 250,
				Tags:            []string{"anime", "rock"},
			},
			wantLen: 3, // similar_bpm, similar_duration, same_tag:anime
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.MatchReasons(tt.seed, tt.candidate)
			if len(got) != tt.wantLen {
				t.Errorf("MatchReasons() returned %d reasons, want %d: %v", len(got), tt.wantLen, got)
			}
		})
	}
}

func TestWeightsForMode(t *testing.T) {
	tests := []struct {
		mode     domain.RecommendMode
		checkBPM float64
		checkTag float64
	}{
		{domain.RecommendModeSimilar, 2.0, 1.0},
		{domain.RecommendModeRelated, 0.5, 3.0},
		{domain.RecommendModeBalanced, 1.5, 2.0},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			weights := WeightsForMode(tt.mode)
			if weights.BPM != tt.checkBPM {
				t.Errorf("BPM weight = %v, want %v", weights.BPM, tt.checkBPM)
			}
			if weights.TagSimilarity != tt.checkTag {
				t.Errorf("TagSimilarity weight = %v, want %v", weights.TagSimilarity, tt.checkTag)
			}
		})
	}
}

func TestSimilarityCalculator_hasGroupRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		artistA *domain.ArtistInfo
		artistB *domain.ArtistInfo
		want    bool
	}{
		{
			name:    "nil artists",
			artistA: nil,
			artistB: nil,
			want:    false,
		},
		{
			name: "A is member of B",
			artistA: &domain.ArtistInfo{
				MBID: "member-123",
				Relations: []domain.MBRelation{
					{Type: "member of band", TargetMBID: "group-456"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "group-456",
			},
			want: true,
		},
		{
			name: "both in same group",
			artistA: &domain.ArtistInfo{
				MBID: "member-1",
				Relations: []domain.MBRelation{
					{Type: "member of band", TargetMBID: "group-456"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "member-2",
				Relations: []domain.MBRelation{
					{Type: "member of band", TargetMBID: "group-456"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.hasGroupRelation(tt.artistA, tt.artistB)
			if got != tt.want {
				t.Errorf("hasGroupRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimilarityCalculator_findCommonTags(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		tagsA   []string
		tagsB   []string
		wantLen int
	}{
		{"no common", []string{"rock"}, []string{"pop"}, 0},
		{"one common", []string{"anime", "rock"}, []string{"anime", "pop"}, 1},
		{"all common", []string{"anime", "jpop"}, []string{"jpop", "anime"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.findCommonTags(tt.tagsA, tt.tagsB)
			if len(got) != tt.wantLen {
				t.Errorf("findCommonTags() returned %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestSimilarityCalculator_hasVoiceActorRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		artistA *domain.ArtistInfo
		artistB *domain.ArtistInfo
		want    bool
	}{
		{
			name:    "nil artists",
			artistA: nil,
			artistB: nil,
			want:    false,
		},
		{
			name: "one nil artist",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
			},
			artistB: nil,
			want:    false,
		},
		{
			name: "shared voice actor by MBID",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Relations: []domain.MBRelation{
					{Type: "voice actor", TargetMBID: "va-123", TargetName: "Voice Actor A"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Relations: []domain.MBRelation{
					{Type: "voice actor", TargetMBID: "va-123", TargetName: "Voice Actor A"},
				},
			},
			want: true,
		},
		{
			name: "shared voice actor by name",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Relations: []domain.MBRelation{
					{Type: "vocal", TargetMBID: "", TargetName: "Singer Name"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Relations: []domain.MBRelation{
					{Type: "vocal", TargetMBID: "", TargetName: "Singer Name"},
				},
			},
			want: true,
		},
		{
			name: "no shared voice actor",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Relations: []domain.MBRelation{
					{Type: "voice actor", TargetMBID: "va-123", TargetName: "Voice Actor A"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Relations: []domain.MBRelation{
					{Type: "voice actor", TargetMBID: "va-456", TargetName: "Voice Actor B"},
				},
			},
			want: false,
		},
		{
			name: "no voice actor relations",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Relations: []domain.MBRelation{
					{Type: "member of band", TargetMBID: "group-123"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Relations: []domain.MBRelation{
					{Type: "member of band", TargetMBID: "group-456"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.hasVoiceActorRelation(tt.artistA, tt.artistB)
			if got != tt.want {
				t.Errorf("hasVoiceActorRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimilarityCalculator_hasCollaborationRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	tests := []struct {
		name    string
		artistA *domain.ArtistInfo
		artistB *domain.ArtistInfo
		want    bool
	}{
		{
			name:    "nil artists",
			artistA: nil,
			artistB: nil,
			want:    false,
		},
		{
			name: "one nil artist",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Name: "Artist 1",
			},
			artistB: nil,
			want:    false,
		},
		{
			name: "A collaborated with B by MBID",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Name: "Artist 1",
				Relations: []domain.MBRelation{
					{Type: "collaboration", TargetMBID: "artist-2", TargetName: "Artist 2"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Name: "Artist 2",
			},
			want: true,
		},
		{
			name: "B collaborated with A by name",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Name: "Artist 1",
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Name: "Artist 2",
				Relations: []domain.MBRelation{
					{Type: "collaborator", TargetMBID: "", TargetName: "Artist 1"},
				},
			},
			want: true,
		},
		{
			name: "no collaboration",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Name: "Artist 1",
				Relations: []domain.MBRelation{
					{Type: "collaboration", TargetMBID: "artist-3", TargetName: "Artist 3"},
				},
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Name: "Artist 2",
				Relations: []domain.MBRelation{
					{Type: "collaboration", TargetMBID: "artist-4", TargetName: "Artist 4"},
				},
			},
			want: false,
		},
		{
			name: "no collaboration relations",
			artistA: &domain.ArtistInfo{
				MBID: "artist-1",
				Name: "Artist 1",
			},
			artistB: &domain.ArtistInfo{
				MBID: "artist-2",
				Name: "Artist 2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.hasCollaborationRelation(tt.artistA, tt.artistB)
			if got != tt.want {
				t.Errorf("hasCollaborationRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimilarityCalculator_durationSimilarity_ZeroDuration(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	// Test with zero durations
	got := calc.durationSimilarity(0, 0)
	if got != 1.0 {
		t.Errorf("durationSimilarity(0, 0) = %v, want 1.0", got)
	}

	// Test with one zero duration - similarity decreases with difference
	got = calc.durationSimilarity(0, 240)
	// Expected: 1.0 - 240/600 = 0.6
	if got < 0.55 || got > 0.65 {
		t.Errorf("durationSimilarity(0, 240) = %v, want ~0.6", got)
	}
}

func TestSimilarityCalculator_gainSimilarity_ZeroGain(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	// Test with zero gains (should use default midpoint)
	got := calc.gainSimilarity(0, 0)
	if got < 0.9 {
		t.Errorf("gainSimilarity(0, 0) = %v, want >= 0.9", got)
	}
}

func TestSimilarityCalculator_hasGroupRelation_ByName(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	// Test group relation matching by name instead of MBID
	artistA := &domain.ArtistInfo{
		MBID: "member-1",
		Relations: []domain.MBRelation{
			{Type: "member of band", TargetMBID: "", TargetName: "Group Name"},
		},
	}
	artistB := &domain.ArtistInfo{
		MBID: "member-2",
		Relations: []domain.MBRelation{
			{Type: "member of", TargetMBID: "", TargetName: "Group Name"},
		},
	}

	got := calc.hasGroupRelation(artistA, artistB)
	if !got {
		t.Error("hasGroupRelation() should return true for members of same group by name")
	}
}

func TestSimilarityCalculator_calculateArtistBonus_GroupRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	artistA := &domain.ArtistInfo{
		MBID:      "member-1",
		SpotifyID: "spotify-1",
		Relations: []domain.MBRelation{
			{Type: "member of band", TargetMBID: "group-123"},
		},
	}
	artistB := &domain.ArtistInfo{
		MBID:      "member-2",
		SpotifyID: "spotify-2",
		Relations: []domain.MBRelation{
			{Type: "member of band", TargetMBID: "group-123"},
		},
	}

	got := calc.calculateArtistBonus(artistA, artistB)
	if got <= 1.0 {
		t.Errorf("calculateArtistBonus() = %v, want > 1.0 for group relation", got)
	}
}

func TestSimilarityCalculator_calculateArtistBonus_VoiceActorRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	artistA := &domain.ArtistInfo{
		MBID:      "artist-1",
		SpotifyID: "spotify-1",
		Relations: []domain.MBRelation{
			{Type: "voice actor", TargetMBID: "va-123"},
		},
	}
	artistB := &domain.ArtistInfo{
		MBID:      "artist-2",
		SpotifyID: "spotify-2",
		Relations: []domain.MBRelation{
			{Type: "voice actor", TargetMBID: "va-123"},
		},
	}

	got := calc.calculateArtistBonus(artistA, artistB)
	if got <= 1.0 {
		t.Errorf("calculateArtistBonus() = %v, want > 1.0 for voice actor relation", got)
	}
}

func TestSimilarityCalculator_calculateArtistBonus_CollaborationRelation(t *testing.T) {
	calc := NewSimilarityCalculator(DefaultWeights(), nil)

	artistA := &domain.ArtistInfo{
		MBID:      "artist-1",
		SpotifyID: "spotify-1",
		Name:      "Artist 1",
		Relations: []domain.MBRelation{
			{Type: "collaboration", TargetMBID: "artist-2", TargetName: "Artist 2"},
		},
	}
	artistB := &domain.ArtistInfo{
		MBID:      "artist-2",
		SpotifyID: "spotify-2",
		Name:      "Artist 2",
	}

	got := calc.calculateArtistBonus(artistA, artistB)
	if got <= 1.0 {
		t.Errorf("calculateArtistBonus() = %v, want > 1.0 for collaboration relation", got)
	}
}
