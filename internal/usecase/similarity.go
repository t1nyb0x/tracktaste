// Package usecase contains business logic for TrackTaste.
package usecase

import (
	"math"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// FeatureWeights defines weights for each audio feature in similarity calculation.
type FeatureWeights struct {
	Tempo        float64
	Energy       float64
	Valence      float64
	Danceability float64
	Acousticness float64
}

// DefaultWeights returns the default feature weights for balanced mode.
func DefaultWeights() FeatureWeights {
	return FeatureWeights{
		Tempo:        1.5,
		Energy:       1.5,
		Valence:      1.2,
		Danceability: 1.0,
		Acousticness: 0.8,
	}
}

// WeightsForMode returns the feature weights for the specified recommendation mode.
func WeightsForMode(mode domain.RecommendMode) FeatureWeights {
	switch mode {
	case domain.RecommendModeSimilar:
		return FeatureWeights{
			Tempo:        2.0,
			Energy:       2.0,
			Valence:      1.5,
			Danceability: 1.2,
			Acousticness: 1.0,
		}
	case domain.RecommendModeRelated:
		return FeatureWeights{
			Tempo:        0.5,
			Energy:       0.5,
			Valence:      0.5,
			Danceability: 0.5,
			Acousticness: 0.5,
		}
	default: // balanced
		return DefaultWeights()
	}
}

// SimilarityCalculator calculates audio feature similarity between tracks.
type SimilarityCalculator struct {
	weights FeatureWeights
}

// NewSimilarityCalculator creates a new SimilarityCalculator with the specified weights.
func NewSimilarityCalculator(weights FeatureWeights) *SimilarityCalculator {
	return &SimilarityCalculator{weights: weights}
}

// Calculate computes the similarity score between two audio features.
// Returns a value between 0.0 and 1.0, where 1.0 means identical.
func (c *SimilarityCalculator) Calculate(seed, candidate *domain.AudioFeatures) float64 {
	if seed == nil || candidate == nil {
		return 0.0
	}

	// Normalize tempo difference (0-250 BPM range)
	tempoDiff := (seed.Tempo - candidate.Tempo) / 250.0

	// Calculate weighted Euclidean distance
	distance := math.Sqrt(
		c.weights.Tempo*tempoDiff*tempoDiff +
			c.weights.Energy*math.Pow(seed.Energy-candidate.Energy, 2) +
			c.weights.Valence*math.Pow(seed.Valence-candidate.Valence, 2) +
			c.weights.Danceability*math.Pow(seed.Danceability-candidate.Danceability, 2) +
			c.weights.Acousticness*math.Pow(seed.Acousticness-candidate.Acousticness, 2),
	)

	// Convert distance to similarity (0.0-1.0)
	similarity := 1.0 / (1.0 + distance)

	return similarity
}

// MatchReasons analyzes which features are similar between two tracks.
// Returns a list of feature names that are considered similar.
func (c *SimilarityCalculator) MatchReasons(seed, candidate *domain.AudioFeatures) []string {
	if seed == nil || candidate == nil {
		return nil
	}

	reasons := make([]string, 0, 5)

	// Threshold for considering features as "similar"
	const threshold = 0.15

	// Tempo: within 15 BPM is considered similar
	if math.Abs(seed.Tempo-candidate.Tempo) <= 15 {
		reasons = append(reasons, "tempo")
	}

	if math.Abs(seed.Energy-candidate.Energy) <= threshold {
		reasons = append(reasons, "energy")
	}

	if math.Abs(seed.Valence-candidate.Valence) <= threshold {
		reasons = append(reasons, "valence")
	}

	if math.Abs(seed.Danceability-candidate.Danceability) <= threshold {
		reasons = append(reasons, "danceability")
	}

	if math.Abs(seed.Acousticness-candidate.Acousticness) <= threshold {
		reasons = append(reasons, "acousticness")
	}

	return reasons
}
