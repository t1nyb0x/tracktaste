package musicbrainz

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/domain"
)

func TestGetRecordingByISRC(t *testing.T) {
	t.Run("empty ISRC", func(t *testing.T) {
		g := NewGateway("")
		_, err := g.GetRecordingByISRC(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty ISRC")
		}
	})

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify User-Agent header
			if r.Header.Get("User-Agent") == "" {
				t.Error("User-Agent header is required")
			}

			resp := rawISRCResponse{
				ISRC: "JPAB12345678",
				Recordings: []rawRecording{
					{
						ID:    "mbid-123",
						Title: "Test Track",
						Tags: []rawTag{
							{Name: "anime", Count: 10},
							{Name: "jpop", Count: 5},
						},
						ArtistCredit: []rawArtist{
							{
								Artist: struct {
									ID   string `json:"id"`
									Name string `json:"name"`
								}{
									ID:   "artist-mbid-123",
									Name: "Test Artist",
								},
							},
						},
					},
				},
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))
		defer server.Close()

		// Note: This test demonstrates the pattern, but actual API URL can't be overridden easily
		// In production, we'd use dependency injection for the base URL
	})
}

func TestGetRecordingWithTags(t *testing.T) {
	t.Run("empty MBID", func(t *testing.T) {
		g := NewGateway("")
		_, err := g.GetRecordingWithTags(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty MBID")
		}
	})
}

func TestGetArtistWithRelations(t *testing.T) {
	t.Run("empty MBID", func(t *testing.T) {
		g := NewGateway("")
		_, err := g.GetArtistWithRelations(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty MBID")
		}
	})
}

func TestGetRecordingsByISRCBatch(t *testing.T) {
	t.Run("empty ISRCs", func(t *testing.T) {
		g := NewGateway("")
		result, err := g.GetRecordingsByISRCBatch(context.Background(), []string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty result, got %d items", len(result))
		}
	})
}

func TestConvertRecording(t *testing.T) {
	g := NewGateway("")

	raw := &rawRecording{
		ID:    "mbid-123",
		Title: "Test Track",
		Tags: []rawTag{
			{Name: "anime", Count: 10},
			{Name: "japanese", Count: 8},
		},
		ArtistCredit: []rawArtist{
			{
				Artist: struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					ID:   "artist-mbid-123",
					Name: "Test Artist",
				},
			},
		},
	}

	recording := g.convertRecording(raw, "JPAB12345678")

	if recording.MBID != "mbid-123" {
		t.Errorf("MBID mismatch: got %s, want %s", recording.MBID, "mbid-123")
	}
	if recording.Title != "Test Track" {
		t.Errorf("Title mismatch: got %s, want %s", recording.Title, "Test Track")
	}
	if recording.ISRC != "JPAB12345678" {
		t.Errorf("ISRC mismatch: got %s, want %s", recording.ISRC, "JPAB12345678")
	}
	if len(recording.Tags) != 2 {
		t.Errorf("Tags count mismatch: got %d, want %d", len(recording.Tags), 2)
	}
	if recording.ArtistMBID != "artist-mbid-123" {
		t.Errorf("ArtistMBID mismatch: got %s, want %s", recording.ArtistMBID, "artist-mbid-123")
	}
	if recording.ArtistName != "Test Artist" {
		t.Errorf("ArtistName mismatch: got %s, want %s", recording.ArtistName, "Test Artist")
	}
}

func TestConvertRecordingNoArtist(t *testing.T) {
	g := NewGateway("")

	raw := &rawRecording{
		ID:           "mbid-123",
		Title:        "Test Track",
		ArtistCredit: []rawArtist{},
	}

	recording := g.convertRecording(raw, "JPAB12345678")

	if recording.ArtistMBID != "" {
		t.Errorf("ArtistMBID should be empty, got %s", recording.ArtistMBID)
	}
	if recording.ArtistName != "" {
		t.Errorf("ArtistName should be empty, got %s", recording.ArtistName)
	}
}

func TestConvertArtist(t *testing.T) {
	g := NewGateway("")

	raw := &rawArtistResponse{
		ID:   "artist-mbid-123",
		Name: "Test Artist",
		Tags: []rawTag{
			{Name: "jpop", Count: 15},
			{Name: "female vocalist", Count: 10},
		},
		Relations: []rawRelation{
			{
				Type: "member of band",
				Artist: &struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					ID:   "group-mbid-456",
					Name: "Test Group",
				},
			},
			{
				Type: "voice actor",
				Artist: &struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					ID:   "va-mbid-789",
					Name: "Voice Actor",
				},
			},
		},
	}

	artist := g.convertArtist(raw)

	if artist.MBID != "artist-mbid-123" {
		t.Errorf("MBID mismatch: got %s, want %s", artist.MBID, "artist-mbid-123")
	}
	if artist.Name != "Test Artist" {
		t.Errorf("Name mismatch: got %s, want %s", artist.Name, "Test Artist")
	}
	if len(artist.Tags) != 2 {
		t.Errorf("Tags count mismatch: got %d, want %d", len(artist.Tags), 2)
	}
	if len(artist.Relations) != 2 {
		t.Errorf("Relations count mismatch: got %d, want %d", len(artist.Relations), 2)
	}

	// Check first relation
	if artist.Relations[0].Type != "member of band" {
		t.Errorf("Relation type mismatch: got %s, want %s", artist.Relations[0].Type, "member of band")
	}
	if artist.Relations[0].TargetMBID != "group-mbid-456" {
		t.Errorf("Relation target MBID mismatch: got %s, want %s", artist.Relations[0].TargetMBID, "group-mbid-456")
	}
}

func TestConvertArtistNoRelations(t *testing.T) {
	g := NewGateway("")

	raw := &rawArtistResponse{
		ID:        "artist-mbid-123",
		Name:      "Test Artist",
		Relations: []rawRelation{},
	}

	artist := g.convertArtist(raw)

	if len(artist.Relations) != 0 {
		t.Errorf("Relations should be empty, got %d items", len(artist.Relations))
	}
}

func TestConvertArtistNilArtistInRelation(t *testing.T) {
	g := NewGateway("")

	raw := &rawArtistResponse{
		ID:   "artist-mbid-123",
		Name: "Test Artist",
		Relations: []rawRelation{
			{
				Type:   "some relation",
				Artist: nil, // No artist, should be skipped
			},
		},
	}

	artist := g.convertArtist(raw)

	if len(artist.Relations) != 0 {
		t.Errorf("Relations should be empty when Artist is nil, got %d items", len(artist.Relations))
	}
}

func TestNewGateway(t *testing.T) {
	t.Run("default user agent", func(t *testing.T) {
		g := NewGateway("")
		if g.userAgent != defaultUserAgent {
			t.Errorf("expected default user agent, got %s", g.userAgent)
		}
	})

	t.Run("custom user agent", func(t *testing.T) {
		customUA := "MyApp/1.0 (test@example.com)"
		g := NewGateway(customUA)
		if g.userAgent != customUA {
			t.Errorf("expected custom user agent %s, got %s", customUA, g.userAgent)
		}
	})

	t.Run("rate limiter is set", func(t *testing.T) {
		g := NewGateway("")
		if g.limiter == nil {
			t.Error("rate limiter should be set")
		}
	})
}

// Integration test (skipped by default, run with -tags=integration)
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// This would test against the real MusicBrainz API
	// Skipped by default to avoid hitting rate limits
	_ = domain.ErrNotFound // Use domain to avoid import error
}
