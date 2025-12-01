// Package v2 contains V2 business logic for TrackTaste.
// V2 uses Deezer + MusicBrainz for track features.
package v2

import (
	"math"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/usecase"
)

// FeatureWeights defines weights for Deezer + MusicBrainz features.
type FeatureWeights struct {
	BPM           float64
	Duration      float64
	Gain          float64
	TagSimilarity float64
}

// DefaultWeights returns the default feature weights for balanced mode.
func DefaultWeights() FeatureWeights {
	return FeatureWeights{
		BPM:           1.5,
		Duration:      0.5,
		Gain:          1.2,
		TagSimilarity: 2.0,
	}
}

// WeightsForMode returns the feature weights for the specified recommendation mode.
func WeightsForMode(mode domain.RecommendMode) FeatureWeights {
	switch mode {
	case domain.RecommendModeSimilar:
		return FeatureWeights{
			BPM:           2.0,
			Duration:      0.8,
			Gain:          1.5,
			TagSimilarity: 1.0,
		}
	case domain.RecommendModeRelated:
		return FeatureWeights{
			BPM:           0.5,
			Duration:      0.3,
			Gain:          0.5,
			TagSimilarity: 3.0,
		}
	default: // balanced
		return DefaultWeights()
	}
}

// SimilarityCalculator calculates similarity using Deezer + MusicBrainz features.
type SimilarityCalculator struct {
	weights      FeatureWeights
	genreMatcher *usecase.GenreMatcher
}

// NewSimilarityCalculator creates a new SimilarityCalculator.
func NewSimilarityCalculator(weights FeatureWeights, genreMatcher *usecase.GenreMatcher) *SimilarityCalculator {
	return &SimilarityCalculator{
		weights:      weights,
		genreMatcher: genreMatcher,
	}
}

// Calculate computes the similarity score between two TrackFeatures.
// Returns a value between 0.0 and 1.0, where 1.0 means identical.
func (c *SimilarityCalculator) Calculate(seed, candidate *domain.TrackFeatures) float64 {
	if seed == nil || candidate == nil {
		return 0.5 // Neutral score when features are unavailable
	}

	var totalWeight float64
	var weightedSum float64

	// BPM similarity
	if seed.BPM > 0 && candidate.BPM > 0 {
		bpmSim := c.bpmSimilarity(seed.BPM, candidate.BPM)
		weightedSum += c.weights.BPM * bpmSim
		totalWeight += c.weights.BPM
	}

	// Duration similarity
	if seed.DurationSeconds > 0 && candidate.DurationSeconds > 0 {
		durSim := c.durationSimilarity(seed.DurationSeconds, candidate.DurationSeconds)
		weightedSum += c.weights.Duration * durSim
		totalWeight += c.weights.Duration
	}

	// Gain similarity
	if seed.Gain != 0 || candidate.Gain != 0 {
		gainSim := c.gainSimilarity(seed.Gain, candidate.Gain)
		weightedSum += c.weights.Gain * gainSim
		totalWeight += c.weights.Gain
	}

	// Tag similarity (Jaccard coefficient)
	if len(seed.Tags) > 0 || len(candidate.Tags) > 0 {
		tagSim := c.tagSimilarity(seed.Tags, candidate.Tags)
		weightedSum += c.weights.TagSimilarity * tagSim
		totalWeight += c.weights.TagSimilarity
	}

	if totalWeight == 0 {
		return 0.5 // Neutral score when no features available
	}

	return weightedSum / totalWeight
}

// CalculateWithBonus computes the final score including genre and artist bonuses.
func (c *SimilarityCalculator) CalculateWithBonus(
	seed, candidate *domain.TrackFeatures,
	seedArtist, candidateArtist *domain.ArtistInfo,
) (baseSimilarity, genreBonus, artistBonus, finalScore float64) {
	baseSimilarity = c.Calculate(seed, candidate)

	// Genre bonus from GenreMatcher
	genreBonus = 1.0
	if c.genreMatcher != nil && seed != nil && candidate != nil {
		genreBonus = c.genreMatcher.CalculateBonus(seed.Tags, candidate.Tags)
	}

	// Artist relation bonus
	artistBonus = c.calculateArtistBonus(seedArtist, candidateArtist)

	finalScore = baseSimilarity * genreBonus * artistBonus
	return
}

// bpmSimilarity calculates BPM similarity.
func (c *SimilarityCalculator) bpmSimilarity(bpmA, bpmB float64) float64 {
	// BPM range: 0-250, normalize the difference
	diff := math.Abs(bpmA-bpmB) / 250.0
	return 1.0 - diff
}

// durationSimilarity calculates duration similarity.
func (c *SimilarityCalculator) durationSimilarity(durA, durB int) float64 {
	// Max duration difference: 10 minutes (600 seconds)
	diff := math.Abs(float64(durA-durB)) / 600.0
	if diff > 1.0 {
		diff = 1.0
	}
	return 1.0 - diff
}

// gainSimilarity calculates gain (loudness) similarity.
func (c *SimilarityCalculator) gainSimilarity(gainA, gainB float64) float64 {
	// Gain range: -20 to 0 dB
	diff := math.Abs(gainA-gainB) / 20.0
	if diff > 1.0 {
		diff = 1.0
	}
	return 1.0 - diff
}

// tagSimilarity calculates Jaccard similarity coefficient for tags.
func (c *SimilarityCalculator) tagSimilarity(tagsA, tagsB []string) float64 {
	if len(tagsA) == 0 && len(tagsB) == 0 {
		return 0.5 // Neutral when both have no tags
	}

	setA := make(map[string]bool)
	for _, tag := range tagsA {
		setA[tag] = true
	}

	setB := make(map[string]bool)
	for _, tag := range tagsB {
		setB[tag] = true
	}

	// Calculate intersection
	intersection := 0
	for tag := range setA {
		if setB[tag] {
			intersection++
		}
	}

	// Calculate union
	union := make(map[string]bool)
	for tag := range setA {
		union[tag] = true
	}
	for tag := range setB {
		union[tag] = true
	}

	if len(union) == 0 {
		return 0.5
	}

	return float64(intersection) / float64(len(union))
}

// calculateArtistBonus calculates bonus based on artist relations.
func (c *SimilarityCalculator) calculateArtistBonus(seedArtist, candidateArtist *domain.ArtistInfo) float64 {
	if seedArtist == nil || candidateArtist == nil {
		return 1.0 // No bonus/penalty
	}

	// Same artist
	if seedArtist.SpotifyID != "" && seedArtist.SpotifyID == candidateArtist.SpotifyID {
		return 1.5
	}
	if seedArtist.MBID != "" && seedArtist.MBID == candidateArtist.MBID {
		return 1.5
	}

	// Check for group/collaboration relations
	if c.hasGroupRelation(seedArtist, candidateArtist) {
		return 1.3
	}

	// Check for voice actor relation (for anime songs)
	if c.hasVoiceActorRelation(seedArtist, candidateArtist) {
		return 1.2
	}

	// Check for collaboration
	if c.hasCollaborationRelation(seedArtist, candidateArtist) {
		return 1.2
	}

	return 1.0
}

// hasGroupRelation checks if artists are in the same group.
func (c *SimilarityCalculator) hasGroupRelation(artistA, artistB *domain.ArtistInfo) bool {
	if artistA == nil || artistB == nil {
		return false
	}

	// Check A's relations for B
	for _, rel := range artistA.Relations {
		if rel.Type == "member of band" || rel.Type == "member of" {
			if rel.TargetMBID == artistB.MBID || rel.TargetName == artistB.Name {
				return true
			}
		}
	}

	// Check B's relations for A
	for _, rel := range artistB.Relations {
		if rel.Type == "member of band" || rel.Type == "member of" {
			if rel.TargetMBID == artistA.MBID || rel.TargetName == artistA.Name {
				return true
			}
		}
	}

	// Check if both are members of the same group
	groupsA := make(map[string]bool)
	for _, rel := range artistA.Relations {
		if rel.Type == "member of band" || rel.Type == "member of" {
			groupsA[rel.TargetMBID] = true
			groupsA[rel.TargetName] = true
		}
	}
	for _, rel := range artistB.Relations {
		if rel.Type == "member of band" || rel.Type == "member of" {
			if groupsA[rel.TargetMBID] || groupsA[rel.TargetName] {
				return true
			}
		}
	}

	return false
}

// hasVoiceActorRelation checks if artists share a voice actor.
func (c *SimilarityCalculator) hasVoiceActorRelation(artistA, artistB *domain.ArtistInfo) bool {
	if artistA == nil || artistB == nil {
		return false
	}

	voiceActorsA := make(map[string]bool)
	for _, rel := range artistA.Relations {
		if rel.Type == "voice actor" || rel.Type == "vocal" {
			voiceActorsA[rel.TargetMBID] = true
			voiceActorsA[rel.TargetName] = true
		}
	}

	for _, rel := range artistB.Relations {
		if rel.Type == "voice actor" || rel.Type == "vocal" {
			if voiceActorsA[rel.TargetMBID] || voiceActorsA[rel.TargetName] {
				return true
			}
		}
	}

	return false
}

// hasCollaborationRelation checks if artists have collaborated.
func (c *SimilarityCalculator) hasCollaborationRelation(artistA, artistB *domain.ArtistInfo) bool {
	if artistA == nil || artistB == nil {
		return false
	}

	for _, rel := range artistA.Relations {
		if rel.Type == "collaboration" || rel.Type == "collaborator" {
			if rel.TargetMBID == artistB.MBID || rel.TargetName == artistB.Name {
				return true
			}
		}
	}

	for _, rel := range artistB.Relations {
		if rel.Type == "collaboration" || rel.Type == "collaborator" {
			if rel.TargetMBID == artistA.MBID || rel.TargetName == artistA.Name {
				return true
			}
		}
	}

	return false
}

// MatchReasons analyzes which features are similar between two tracks.
func (c *SimilarityCalculator) MatchReasons(seed, candidate *domain.TrackFeatures) []string {
	if seed == nil || candidate == nil {
		return nil
	}

	reasons := make([]string, 0, 5)

	// BPM: within 15 BPM is considered similar
	if seed.BPM > 0 && candidate.BPM > 0 {
		if math.Abs(seed.BPM-candidate.BPM) <= 15 {
			reasons = append(reasons, "similar_bpm")
		}
	}

	// Duration: within 30 seconds is considered similar
	if seed.DurationSeconds > 0 && candidate.DurationSeconds > 0 {
		if math.Abs(float64(seed.DurationSeconds-candidate.DurationSeconds)) <= 30 {
			reasons = append(reasons, "similar_duration")
		}
	}

	// Gain: within 3 dB is considered similar
	if seed.Gain != 0 || candidate.Gain != 0 {
		if math.Abs(seed.Gain-candidate.Gain) <= 3 {
			reasons = append(reasons, "similar_loudness")
		}
	}

	// Common tags
	commonTags := c.findCommonTags(seed.Tags, candidate.Tags)
	for _, tag := range commonTags {
		reasons = append(reasons, "same_tag:"+tag)
	}

	return reasons
}

// findCommonTags returns common tags between two tag lists.
func (c *SimilarityCalculator) findCommonTags(tagsA, tagsB []string) []string {
	setA := make(map[string]bool)
	for _, tag := range tagsA {
		setA[tag] = true
	}

	common := make([]string, 0)
	for _, tag := range tagsB {
		if setA[tag] {
			common = append(common, tag)
		}
	}

	return common
}
