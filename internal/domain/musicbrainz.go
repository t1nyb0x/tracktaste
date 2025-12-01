// Package domain defines the core business entities for TrackTaste.
package domain

// MBRecording represents a recording from MusicBrainz.
type MBRecording struct {
	MBID       string  `json:"mbid"`
	Title      string  `json:"title"`
	ISRC       string  `json:"isrc"`
	Tags       []MBTag `json:"tags"`
	ArtistMBID string  `json:"artist_mbid"`
	ArtistName string  `json:"artist_name"`
}

// MBTag represents a user-contributed tag from MusicBrainz.
type MBTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"` // Vote count (can be used for weighting)
}

// MBArtist represents an artist from MusicBrainz.
type MBArtist struct {
	MBID      string       `json:"mbid"`
	Name      string       `json:"name"`
	Tags      []MBTag      `json:"tags"`
	Relations []MBRelation `json:"relations"`
}

// MBRelation represents a relation between artists.
type MBRelation struct {
	Type       string `json:"type"`        // e.g., "member of band", "voice actor", "collaboration"
	TargetMBID string `json:"target_mbid"` // MBID of the related artist
	TargetName string `json:"target_name"` // Name of the related artist
}
