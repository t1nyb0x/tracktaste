// Package domain defines the core business entities for TrackTaste.
package domain

// AudioFeatures represents the audio features of a track.
// These features are retrieved from Spotify Audio Features API.
type AudioFeatures struct {
	TrackID          string  `json:"track_id"`
	Tempo            float64 `json:"tempo"`            // BPM (0-250)
	Energy           float64 `json:"energy"`           // 0.0-1.0
	Danceability     float64 `json:"danceability"`     // 0.0-1.0
	Valence          float64 `json:"valence"`          // 0.0-1.0 (positiveness)
	Acousticness     float64 `json:"acousticness"`     // 0.0-1.0
	Instrumentalness float64 `json:"instrumentalness"` // 0.0-1.0
	Speechiness      float64 `json:"speechiness"`      // 0.0-1.0
	Liveness         float64 `json:"liveness"`         // 0.0-1.0
	Loudness         float64 `json:"loudness"`         // dB (-60 to 0)
	Key              int     `json:"key"`              // 0-11 (C=0, C#=1, ...)
	Mode             int     `json:"mode"`             // 0=minor, 1=major
	TimeSignature    int     `json:"time_signature"`   // 3-7 (4 = 4/4)
	DurationMs       int     `json:"duration_ms"`
}
