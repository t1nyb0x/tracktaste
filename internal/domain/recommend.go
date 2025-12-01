// Package domain defines the core business entities for TrackTaste.
package domain

// RecommendMode represents the recommendation mode.
type RecommendMode string

const (
	// RecommendModeSimilar prioritizes audio features similarity.
	RecommendModeSimilar RecommendMode = "similar"
	// RecommendModeRelated prioritizes artist/genre relationships.
	RecommendModeRelated RecommendMode = "related"
	// RecommendModeBalanced balances both audio features and relationships.
	RecommendModeBalanced RecommendMode = "balanced"
)

// ParseRecommendMode parses a string into RecommendMode.
// Returns RecommendModeBalanced if the input is invalid.
func ParseRecommendMode(s string) RecommendMode {
	switch s {
	case "similar":
		return RecommendModeSimilar
	case "related":
		return RecommendModeRelated
	case "balanced":
		return RecommendModeBalanced
	default:
		return RecommendModeBalanced
	}
}

// RecommendedTrack represents a recommended track with similarity information.
type RecommendedTrack struct {
	Track           Track          `json:"track"`
	SimilarityScore float64        `json:"similarity_score"`
	GenreBonus      float64        `json:"genre_bonus"`
	FinalScore      float64        `json:"final_score"`
	MatchReasons    []string       `json:"match_reasons"`
	Features        *TrackFeatures `json:"features,omitempty"`
	// Deprecated: Use Features instead
	AudioFeatures *AudioFeatures `json:"audio_features,omitempty"`
}

// RecommendResult represents the result of a recommendation request.
type RecommendResult struct {
	SeedTrack    Track              `json:"seed_track"`
	SeedFeatures *TrackFeatures     `json:"seed_features,omitempty"`
	SeedGenres   []string           `json:"seed_genres,omitempty"`
	Items        []RecommendedTrack `json:"items"`
	Mode         RecommendMode      `json:"mode"`
	// Deprecated: Use SeedFeatures instead
	SeedAudioFeatures *AudioFeatures `json:"seed_audio_features,omitempty"`
}
