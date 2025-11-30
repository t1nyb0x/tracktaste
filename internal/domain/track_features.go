// Package domain defines the core business entities for TrackTaste.
package domain

// TrackFeatures represents the combined audio features from Deezer and MusicBrainz.
// This replaces the deprecated Spotify AudioFeatures.
type TrackFeatures struct {
	TrackID string `json:"track_id"`
	ISRC    string `json:"isrc"`

	// Deezer features
	BPM             float64 `json:"bpm"`              // Tempo (0-250)
	DurationSeconds int     `json:"duration_seconds"` // Track duration in seconds
	Gain            float64 `json:"gain"`             // ReplayGain (dB)

	// MusicBrainz features
	Tags       []string `json:"tags"`        // Genre/style tags
	ArtistMBID string   `json:"artist_mbid"` // MusicBrainz Artist ID
}

// HasDeezerFeatures returns true if Deezer features are available.
func (f *TrackFeatures) HasDeezerFeatures() bool {
	return f.BPM > 0 || f.DurationSeconds > 0 || f.Gain != 0
}

// HasMusicBrainzFeatures returns true if MusicBrainz features are available.
func (f *TrackFeatures) HasMusicBrainzFeatures() bool {
	return len(f.Tags) > 0 || f.ArtistMBID != ""
}

// GetTagNames returns a list of tag names for comparison.
func (f *TrackFeatures) GetTagNames() []string {
	if f.Tags == nil {
		return []string{}
	}
	return f.Tags
}

// ArtistInfo represents artist information for relation bonus calculation.
type ArtistInfo struct {
	SpotifyID  string       // Spotify Artist ID
	MBID       string       // MusicBrainz Artist ID
	Name       string       // Artist name
	Genres     []string     // Spotify genres
	Tags       []MBTag      // MusicBrainz tags
	Relations  []MBRelation // MusicBrainz artist relations
}
