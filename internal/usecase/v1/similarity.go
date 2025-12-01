// Package v1 contains V1 business logic for TrackTaste.
package v1

import (
	"math"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

// FeatureWeights defines the weights for Spotify Audio Features.
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
		Valence:      1.5,
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
			Danceability: 1.0,
			Acousticness: 1.0,
		}
	case domain.RecommendModeRelated:
		return FeatureWeights{
			Tempo:        0.5,
			Energy:       0.8,
			Valence:      1.0,
			Danceability: 0.5,
			Acousticness: 0.5,
		}
	default: // balanced
		return DefaultWeights()
	}
}

// SimilarityCalculator calculates similarity between tracks based on Spotify Audio Features.
type SimilarityCalculator struct {
	weights FeatureWeights
}

// NewSimilarityCalculator creates a new SimilarityCalculator with the given weights.
func NewSimilarityCalculator(weights FeatureWeights) *SimilarityCalculator {
	return &SimilarityCalculator{weights: weights}
}

// Calculate computes the similarity score between two audio features.
// Returns a value between 0.0 and 1.0, where 1.0 means identical.
func (c *SimilarityCalculator) Calculate(seed, candidate *domain.AudioFeatures) float64 {
	if seed == nil || candidate == nil {
		return 0.0
	}

	var totalWeight float64
	var weightedSum float64

	// Tempo similarity (normalize to 0-200 BPM range)
	tempoSim := c.tempoSimilarity(seed.Tempo, candidate.Tempo)
	weightedSum += c.weights.Tempo * tempoSim
	totalWeight += c.weights.Tempo

	// Energy similarity (0-1 range)
	energySim := 1 - math.Abs(seed.Energy-candidate.Energy)
	weightedSum += c.weights.Energy * energySim
	totalWeight += c.weights.Energy

	// Valence similarity (0-1 range)
	valenceSim := 1 - math.Abs(seed.Valence-candidate.Valence)
	weightedSum += c.weights.Valence * valenceSim
	totalWeight += c.weights.Valence

	// Danceability similarity (0-1 range)
	danceabilitySim := 1 - math.Abs(seed.Danceability-candidate.Danceability)
	weightedSum += c.weights.Danceability * danceabilitySim
	totalWeight += c.weights.Danceability

	// Acousticness similarity (0-1 range)
	acousticnessSim := 1 - math.Abs(seed.Acousticness-candidate.Acousticness)
	weightedSum += c.weights.Acousticness * acousticnessSim
	totalWeight += c.weights.Acousticness

	if totalWeight == 0 {
		return 0.0
	}

	return weightedSum / totalWeight
}

// tempoSimilarity calculates tempo similarity accounting for tempo doubling/halving.
func (c *SimilarityCalculator) tempoSimilarity(tempoA, tempoB float64) float64 {
	if tempoA <= 0 || tempoB <= 0 {
		return 0.5 // Neutral score when tempo is unknown
	}

	// Consider tempo doubling/halving (60 BPM ≈ 120 BPM)
	ratios := []float64{
		tempoA / tempoB,
		(tempoA * 2) / tempoB,
		tempoA / (tempoB * 2),
	}

	minDiff := math.MaxFloat64
	for _, ratio := range ratios {
		diff := math.Abs(1 - ratio)
		if diff < minDiff {
			minDiff = diff
		}
	}

	// Normalize: diff of 0.5 (50% tempo difference) = 0 similarity
	similarity := 1 - (minDiff * 2)
	if similarity < 0 {
		similarity = 0
	}

	return similarity
}

// MatchReasons returns a list of features that are similar between two tracks.
func (c *SimilarityCalculator) MatchReasons(seed, candidate *domain.AudioFeatures) []string {
	if seed == nil || candidate == nil {
		return nil
	}

	var reasons []string
	threshold := 0.15 // 15% difference threshold

	// Check tempo
	tempoSim := c.tempoSimilarity(seed.Tempo, candidate.Tempo)
	if tempoSim > 0.85 {
		reasons = append(reasons, "tempo")
	}

	// Check energy
	if math.Abs(seed.Energy-candidate.Energy) < threshold {
		reasons = append(reasons, "energy")
	}

	// Check valence
	if math.Abs(seed.Valence-candidate.Valence) < threshold {
		reasons = append(reasons, "valence")
	}

	// Check danceability
	if math.Abs(seed.Danceability-candidate.Danceability) < threshold {
		reasons = append(reasons, "danceability")
	}

	// Check acousticness
	if math.Abs(seed.Acousticness-candidate.Acousticness) < threshold {
		reasons = append(reasons, "acousticness")
	}

	return reasons
}
