package usecase

import "testing"

func TestGenreMatcher_CalculateBonus(t *testing.T) {
	tests := []struct {
		name            string
		seedGenres      []string
		candidateGenres []string
		want            float64
	}{
		{
			name:            "exact match",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"anime"},
			want:            2.0, // Strong boost for exact match
		},
		{
			name:            "exact match with multiple genres",
			seedGenres:      []string{"j-pop", "anime"},
			candidateGenres: []string{"k-pop", "anime"},
			want:            2.0, // Strong boost for exact match
		},
		{
			name:            "same group (otaku)",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"japanese vgm"},
			want:            1.5, // Same group bonus
		},
		{
			name:            "same group (jpop)",
			seedGenres:      []string{"j-pop"},
			candidateGenres: []string{"city pop"},
			want:            1.5, // Same group bonus
		},
		{
			name:            "related groups (otaku <-> rock)",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"j-rock"},
			want:            1.0, // Related groups = neutral
		},
		{
			name:            "related groups (idol <-> jpop)",
			seedGenres:      []string{"japanese idol"},
			candidateGenres: []string{"j-pop"},
			want:            1.0, // Related groups = neutral
		},
		{
			name:            "unrelated groups (otaku <-> kpop)",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"k-pop"},
			want:            0.3, // Strong penalty for unrelated
		},
		{
			name:            "unrelated groups (otaku <-> jpop)",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"j-pop"},
			want:            0.3, // J-POP is now unrelated to Otaku
		},
		{
			name:            "empty seed genres",
			seedGenres:      []string{},
			candidateGenres: []string{"anime"},
			want:            1.0,
		},
		{
			name:            "empty candidate genres",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{},
			want:            1.0,
		},
		{
			name:            "unknown genres both",
			seedGenres:      []string{"unknown-genre"},
			candidateGenres: []string{"another-unknown"},
			want:            1.0,
		},
		{
			name:            "unknown seed vs known candidate",
			seedGenres:      []string{"unknown-genre"},
			candidateGenres: []string{"anime"},
			want:            0.3, // Strong penalty for mismatch
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewGenreMatcher()
			got := m.CalculateBonus(tt.seedGenres, tt.candidateGenres)
			if got != tt.want {
				t.Errorf("CalculateBonus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenreMatcher_IsGenreMatch(t *testing.T) {
	tests := []struct {
		name            string
		seedGenres      []string
		candidateGenres []string
		want            bool
	}{
		{
			name:            "exact match",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"anime"},
			want:            true,
		},
		{
			name:            "same group",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"japanese vgm"},
			want:            true,
		},
		{
			name:            "different groups",
			seedGenres:      []string{"anime"},
			candidateGenres: []string{"k-pop"},
			want:            false,
		},
		{
			name:            "unknown genres - different strings",
			seedGenres:      []string{"some-unknown-genre"},
			candidateGenres: []string{"another-unknown-genre"},
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewGenreMatcher()
			got := m.IsGenreMatch(tt.seedGenres, tt.candidateGenres)
			if got != tt.want {
				t.Errorf("IsGenreMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGenreGroup(t *testing.T) {
	tests := []struct {
		genres []string
		want   GenreGroup
	}{
		{[]string{"anime"}, GenreGroupOtaku},
		{[]string{"japanese vgm"}, GenreGroupOtaku},
		{[]string{"vocaloid"}, GenreGroupOtaku},
		{[]string{"j-pop"}, GenreGroupJPop},
		{[]string{"city pop"}, GenreGroupJPop},
		{[]string{"j-rock"}, GenreGroupRock},
		{[]string{"visual kei"}, GenreGroupRock},
		{[]string{"k-pop"}, GenreGroupKPop},
		{[]string{"japanese idol"}, GenreGroupIdol},
		{[]string{"unknown"}, GenreGroupOther},
		{[]string{}, GenreGroupOther},
		// Priority: otaku > idol > jpop > rock > kpop
		{[]string{"anime", "j-pop"}, GenreGroupOtaku},
		{[]string{"japanese idol", "j-rock"}, GenreGroupIdol},
	}

	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			got := getGenreGroup(tt.genres)
			if got != tt.want {
				t.Errorf("getGenreGroup(%v) = %v, want %v", tt.genres, got, tt.want)
			}
		})
	}
}

func (g GenreGroup) String() string {
	return string(g)
}
