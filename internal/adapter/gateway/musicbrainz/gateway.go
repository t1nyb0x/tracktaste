// Package musicbrainz provides the MusicBrainz API gateway implementation.
package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
	"golang.org/x/time/rate"
)

const (
	apiBaseURL = "https://musicbrainz.org/ws/2"
	// MusicBrainz rate limit: 1 request per second
	// User-Agent is required
	defaultUserAgent = "TrackTaste/1.0 (https://github.com/t1nyb0x/tracktaste)"
)

// Gateway implements the MusicBrainzAPI interface.
type Gateway struct {
	httpc     *http.Client
	limiter   *rate.Limiter
	userAgent string
}

// NewGateway creates a new MusicBrainz API gateway.
// userAgent should be in the format "AppName/Version (contact-url-or-email)"
func NewGateway(userAgent string) *Gateway {
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	return &Gateway{
		httpc:     &http.Client{Timeout: 30 * time.Second},
		limiter:   rate.NewLimiter(rate.Every(time.Second), 1), // 1 req/sec
		userAgent: userAgent,
	}
}

// rawISRCResponse represents the response from ISRC lookup.
type rawISRCResponse struct {
	ISRC       string         `json:"isrc"`
	Recordings []rawRecording `json:"recordings"`
}

// rawRecording represents a recording from MusicBrainz.
type rawRecording struct {
	ID            string      `json:"id"` // MBID
	Title         string      `json:"title"`
	ArtistCredit  []rawArtist `json:"artist-credit"`
	Tags          []rawTag    `json:"tags"`
	ISRCs         []string    `json:"isrcs"`
	Disambiguation string     `json:"disambiguation"`
}

// rawArtist represents an artist in MusicBrainz response.
type rawArtist struct {
	Artist struct {
		ID   string `json:"id"` // MBID
		Name string `json:"name"`
	} `json:"artist"`
	Name string `json:"name"` // credited name
}

// rawTag represents a tag in MusicBrainz response.
type rawTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// rawArtistResponse represents a full artist response.
type rawArtistResponse struct {
	ID        string        `json:"id"` // MBID
	Name      string        `json:"name"`
	Tags      []rawTag      `json:"tags"`
	Relations []rawRelation `json:"relations"`
}

// rawRelation represents a relation between entities.
type rawRelation struct {
	Type   string `json:"type"`
	Artist *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
}

// GetRecordingByISRC searches for a recording by ISRC.
func (g *Gateway) GetRecordingByISRC(ctx context.Context, isrc string) (*domain.MBRecording, error) {
	if isrc == "" {
		return nil, fmt.Errorf("musicbrainz: ISRC is required")
	}

	// Wait for rate limiter
	if err := g.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("musicbrainz: rate limiter error: %w", err)
	}

	endpoint := fmt.Sprintf("%s/isrc/%s?inc=recordings+artists+tags&fmt=json", apiBaseURL, isrc)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", g.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: unexpected status %d", resp.StatusCode)
	}

	var raw rawISRCResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to decode response: %w", err)
	}

	if len(raw.Recordings) == 0 {
		return nil, domain.ErrNotFound
	}

	// Return the first recording
	return g.convertRecording(&raw.Recordings[0], isrc), nil
}

// GetRecordingWithTags retrieves recording details including tags.
func (g *Gateway) GetRecordingWithTags(ctx context.Context, mbid string) (*domain.MBRecording, error) {
	if mbid == "" {
		return nil, fmt.Errorf("musicbrainz: MBID is required")
	}

	// Wait for rate limiter
	if err := g.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("musicbrainz: rate limiter error: %w", err)
	}

	endpoint := fmt.Sprintf("%s/recording/%s?inc=tags+artists&fmt=json", apiBaseURL, mbid)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", g.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: unexpected status %d", resp.StatusCode)
	}

	var raw rawRecording
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to decode response: %w", err)
	}

	return g.convertRecording(&raw, ""), nil
}

// GetArtistWithRelations retrieves artist details including tags and relations.
func (g *Gateway) GetArtistWithRelations(ctx context.Context, mbid string) (*domain.MBArtist, error) {
	if mbid == "" {
		return nil, fmt.Errorf("musicbrainz: MBID is required")
	}

	// Wait for rate limiter
	if err := g.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("musicbrainz: rate limiter error: %w", err)
	}

	endpoint := fmt.Sprintf("%s/artist/%s?inc=tags+artist-rels&fmt=json", apiBaseURL, mbid)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", g.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: unexpected status %d", resp.StatusCode)
	}

	var raw rawArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to decode response: %w", err)
	}

	return g.convertArtist(&raw), nil
}

// GetRecordingsByISRCBatch retrieves multiple recordings by their ISRCs.
// MusicBrainz does not have a batch endpoint, so we make sequential requests
// due to the strict rate limit.
func (g *Gateway) GetRecordingsByISRCBatch(ctx context.Context, isrcs []string) (map[string]*domain.MBRecording, error) {
	if len(isrcs) == 0 {
		return make(map[string]*domain.MBRecording), nil
	}

	result := make(map[string]*domain.MBRecording)
	var mu sync.Mutex

	// MusicBrainz has a strict rate limit of 1 req/sec, so we process sequentially
	// but we can parallelize up to the rate limit with proper spacing
	for _, isrc := range isrcs {
		recording, err := g.GetRecordingByISRC(ctx, isrc)
		if err != nil {
			if err != domain.ErrNotFound {
				logger.Warning("MusicBrainz", fmt.Sprintf("Failed to get recording by ISRC %s: %v", isrc, err))
			}
			continue
		}

		mu.Lock()
		result[isrc] = recording
		mu.Unlock()
	}

	return result, nil
}

// convertRecording converts raw MusicBrainz recording to domain model.
func (g *Gateway) convertRecording(raw *rawRecording, isrc string) *domain.MBRecording {
	recording := &domain.MBRecording{
		MBID:  raw.ID,
		Title: raw.Title,
		ISRC:  isrc,
		Tags:  make([]domain.MBTag, 0, len(raw.Tags)),
	}

	// Convert tags
	for _, tag := range raw.Tags {
		recording.Tags = append(recording.Tags, domain.MBTag{
			Name:  tag.Name,
			Count: tag.Count,
		})
	}

	// Get primary artist info
	if len(raw.ArtistCredit) > 0 {
		recording.ArtistMBID = raw.ArtistCredit[0].Artist.ID
		recording.ArtistName = raw.ArtistCredit[0].Artist.Name
	}

	return recording
}

// convertArtist converts raw MusicBrainz artist to domain model.
func (g *Gateway) convertArtist(raw *rawArtistResponse) *domain.MBArtist {
	artist := &domain.MBArtist{
		MBID:      raw.ID,
		Name:      raw.Name,
		Tags:      make([]domain.MBTag, 0, len(raw.Tags)),
		Relations: make([]domain.MBRelation, 0),
	}

	// Convert tags
	for _, tag := range raw.Tags {
		artist.Tags = append(artist.Tags, domain.MBTag{
			Name:  tag.Name,
			Count: tag.Count,
		})
	}

	// Convert relations (only artist-to-artist relations)
	for _, rel := range raw.Relations {
		if rel.Artist != nil {
			artist.Relations = append(artist.Relations, domain.MBRelation{
				Type:       rel.Type,
				TargetMBID: rel.Artist.ID,
				TargetName: rel.Artist.Name,
			})
		}
	}

	return artist
}

// rawBrowseRecordingsResponse represents the response from browse recordings.
type rawBrowseRecordingsResponse struct {
	RecordingCount int            `json:"recording-count"`
	RecordingOffset int           `json:"recording-offset"`
	Recordings     []rawRecording `json:"recordings"`
}

// GetArtistRecordings retrieves recordings by an artist (same artist's other tracks).
func (g *Gateway) GetArtistRecordings(ctx context.Context, artistMBID string, limit int) ([]domain.MBRecording, error) {
	if artistMBID == "" {
		return nil, fmt.Errorf("musicbrainz: artist MBID is required")
	}

	if limit <= 0 {
		limit = 25 // MusicBrainz default
	}
	if limit > 100 {
		limit = 100 // MusicBrainz max
	}

	// Wait for rate limiter
	if err := g.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("musicbrainz: rate limiter error: %w", err)
	}

	// Browse recordings by artist, include ISRCs
	endpoint := fmt.Sprintf("%s/recording?artist=%s&inc=isrcs+tags&limit=%d&fmt=json",
		apiBaseURL, artistMBID, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", g.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("musicbrainz: unexpected status %d", resp.StatusCode)
	}

	var raw rawBrowseRecordingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("musicbrainz: failed to decode response: %w", err)
	}

	result := make([]domain.MBRecording, 0, len(raw.Recordings))
	for _, rec := range raw.Recordings {
		isrc := ""
		if len(rec.ISRCs) > 0 {
			isrc = rec.ISRCs[0]
		}
		converted := g.convertRecording(&rec, isrc)
		converted.ArtistMBID = artistMBID // Set artist MBID since we know it
		result = append(result, *converted)
	}

	logger.Debug("MusicBrainz", fmt.Sprintf("Got %d recordings for artist %s", len(result), artistMBID))
	return result, nil
}
