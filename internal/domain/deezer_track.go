// Package domain defines the core business entities for TrackTaste.
package domain

// DeezerTrack represents track information retrieved from Deezer API.
// It contains BPM, Duration, and Gain information used for similarity calculation.
type DeezerTrack struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	ISRC            string  `json:"isrc"`
	BPM             float64 `json:"bpm"`              // Tempo (beats per minute)
	DurationSeconds int     `json:"duration_seconds"` // Track duration in seconds
	Gain            float64 `json:"gain"`             // ReplayGain value in dB
	ExplicitLyrics  bool    `json:"explicit_lyrics"`
	ArtistID        int64   `json:"artist_id"`
	ArtistName      string  `json:"artist_name"`
}
