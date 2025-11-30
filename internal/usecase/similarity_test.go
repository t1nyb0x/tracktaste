package usecase

import (
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

func TestSimilarityCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name      string
		seed      *domain.AudioFeatures
		candidate *domain.AudioFeatures
		weights   FeatureWeights
		wantMin   float64
		wantMax   float64
	}{
		{
			name: "identical features should return high similarity",
			seed: &domain.AudioFeatures{
				Tempo:        128.0,
				Energy:       0.8,
				Valence:      0.6,
				Danceability: 0.7,
				Acousticness: 0.2,
			},
			candidate: &domain.AudioFeatures{
				Tempo:        128.0,
				Energy:       0.8,
				Valence:      0.6,
				Danceability: 0.7,
				Acousticness: 0.2,
			},
			weights: DefaultWeights(),
			wantMin: 0.99,
			wantMax: 1.0,
		},
		{
			name: "similar features should return high similarity",
			seed: &domain.AudioFeatures{
				Tempo:        128.0,
				Energy:       0.8,
				Valence:      0.6,
				Danceability: 0.7,
				Acousticness: 0.2,
			},
			candidate: &domain.AudioFeatures{
				Tempo:        130.0,
				Energy:       0.78,
				Valence:      0.62,
				Danceability: 0.68,
				Acousticness: 0.22,
			},
			weights: DefaultWeights(),
			wantMin: 0.9,
			wantMax: 1.0,
		},
		{
			name: "very different features should return low similarity",
			seed: &domain.AudioFeatures{
				Tempo:        60.0,
				Energy:       0.1,
				Valence:      0.2,
				Danceability: 0.3,
				Acousticness: 0.9,
			},
			candidate: &domain.AudioFeatures{
				Tempo:        180.0,
				Energy:       0.9,
				Valence:      0.9,
				Danceability: 0.9,
				Acousticness: 0.1,
			},
			weights: DefaultWeights(),
			wantMin: 0.0,
			wantMax: 0.5,
		},
		{
			name: "nil seed should return 0",
			seed: nil,
			candidate: &domain.AudioFeatures{
				Tempo:  128.0,
				Energy: 0.8,
			},
			weights: DefaultWeights(),
			wantMin: 0.0,
			wantMax: 0.0,
		},
		{
			name: "nil candidate should return 0",
			seed: &domain.AudioFeatures{
				Tempo:  128.0,
				Energy: 0.8,
			},
			candidate: nil,
			weights:   DefaultWeights(),
			wantMin:   0.0,
			wantMax:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSimilarityCalculator(tt.weights)
			got := c.Calculate(tt.seed, tt.candidate)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("Calculate() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSimilarityCalculator_MatchReasons(t *testing.T) {
	tests := []struct {
		name      string
		seed      *domain.AudioFeatures
		candidate *domain.AudioFeatures
		want      []string
	}{
		{
			name: "all features similar",
			seed: &domain.AudioFeatures{
				Tempo:        128.0,
				Energy:       0.8,
				Valence:      0.6,
				Danceability: 0.7,
				Acousticness: 0.2,
			},
			candidate: &domain.AudioFeatures{
				Tempo:        130.0,
				Energy:       0.82,
				Valence:      0.62,
				Danceability: 0.72,
				Acousticness: 0.22,
			},
			want: []string{"tempo", "energy", "valence", "danceability", "acousticness"},
		},
		{
			name: "only tempo similar",
			seed: &domain.AudioFeatures{
				Tempo:        128.0,
				Energy:       0.1,
				Valence:      0.1,
				Danceability: 0.1,
				Acousticness: 0.1,
			},
			candidate: &domain.AudioFeatures{
				Tempo:        130.0,
				Energy:       0.9,
				Valence:      0.9,
				Danceability: 0.9,
				Acousticness: 0.9,
			},
			want: []string{"tempo"},
		},
		{
			name: "nil features",
			seed: nil,
			candidate: &domain.AudioFeatures{
				Tempo: 128.0,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewSimilarityCalculator(DefaultWeights())
			got := c.MatchReasons(tt.seed, tt.candidate)

			if tt.want == nil {
				if got != nil {
					t.Errorf("MatchReasons() = %v, want nil", got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("MatchReasons() = %v, want %v", got, tt.want)
				return
			}

			gotMap := make(map[string]bool)
			for _, r := range got {
				gotMap[r] = true
			}
			for _, w := range tt.want {
				if !gotMap[w] {
					t.Errorf("MatchReasons() missing %v, got %v", w, got)
				}
			}
		})
	}
}

func TestWeightsForMode(t *testing.T) {
	tests := []struct {
		mode      domain.RecommendMode
		wantTempo float64
	}{
		{domain.RecommendModeSimilar, 2.0},
		{domain.RecommendModeRelated, 0.5},
		{domain.RecommendModeBalanced, 1.5},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			weights := WeightsForMode(tt.mode)
			if weights.Tempo != tt.wantTempo {
				t.Errorf("WeightsForMode(%s).Tempo = %v, want %v", tt.mode, weights.Tempo, tt.wantTempo)
			}
		})
	}
}
