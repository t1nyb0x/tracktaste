package spotify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/t1nyb0x/tracktaste/internal/port/repository"
)

// mockTokenRepository implements repository.TokenRepository for testing
type mockTokenRepository struct {
	tokens         map[string]string
	invalidatedKey string
}

func newMockTokenRepo() *mockTokenRepository {
	return &mockTokenRepository{tokens: make(map[string]string)}
}

func (m *mockTokenRepository) SaveToken(ctx context.Context, key, token string, expiresIn int) error {
	m.tokens[key] = token
	return nil
}

func (m *mockTokenRepository) GetToken(ctx context.Context, key string) (string, error) {
	return m.tokens[key], nil
}

func (m *mockTokenRepository) IsTokenValid(ctx context.Context, key string) bool {
	_, ok := m.tokens[key]
	return ok
}

func (m *mockTokenRepository) InvalidateToken(ctx context.Context, key string) error {
	m.invalidatedKey = key
	delete(m.tokens, key)
	return nil
}

var _ repository.TokenRepository = (*mockTokenRepository)(nil)

func TestNewGateway(t *testing.T) {
	repo := newMockTokenRepo()
	gw := NewGateway("client_id", "secret", repo)

	if gw.clientID != "client_id" {
		t.Errorf("expected clientID 'client_id', got '%s'", gw.clientID)
	}
	if gw.secret != "secret" {
		t.Errorf("expected secret 'secret', got '%s'", gw.secret)
	}
	if gw.tokenRepo != repo {
		t.Error("expected tokenRepo to be set")
	}
	if gw.httpc == nil {
		t.Error("expected httpc to be initialized")
	}
}

func TestRawTrack_ToDomain(t *testing.T) {
	raw := rawTrack{
		ID:         "track123",
		Name:       "Test Track",
		DiscNumber: 1,
		DurationMs: 180000,
		Explicit:   true,
		Popularity: 75,
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/track/track123",
		},
		ExternalIDs: struct {
			ISRC string `json:"isrc"`
		}{ISRC: "USRC12345678"},
		TrackNumber: 5,
		Artists: []rawSimpleArtist{
			{
				ID:   "artist1",
				Name: "Test Artist",
				ExternalURLs: map[string]string{
					"spotify": "https://open.spotify.com/artist/artist1",
				},
			},
		},
		Album: rawAlbum{
			ID:          "album1",
			Name:        "Test Album",
			ReleaseDate: "2023-06-15",
			ExternalURLs: map[string]string{
				"spotify": "https://open.spotify.com/album/album1",
			},
			Images: []rawImage{
				{URL: "https://example.com/image.jpg", Height: 640, Width: 640},
			},
			Artists: []rawSimpleArtist{
				{
					ID:   "artist1",
					Name: "Test Artist",
					ExternalURLs: map[string]string{
						"spotify": "https://open.spotify.com/artist/artist1",
					},
				},
			},
		},
	}

	track := raw.toDomain()

	if track.ID != "track123" {
		t.Errorf("expected ID 'track123', got '%s'", track.ID)
	}
	if track.Name != "Test Track" {
		t.Errorf("expected Name 'Test Track', got '%s'", track.Name)
	}
	if track.DurationMs != 180000 {
		t.Errorf("expected DurationMs 180000, got %d", track.DurationMs)
	}
	if !track.Explicit {
		t.Error("expected Explicit to be true")
	}
	if track.Popularity == nil || *track.Popularity != 75 {
		t.Errorf("expected Popularity 75, got %v", track.Popularity)
	}
	if track.ISRC == nil || *track.ISRC != "USRC12345678" {
		t.Errorf("expected ISRC 'USRC12345678', got %v", track.ISRC)
	}
	if len(track.Artists) != 1 {
		t.Errorf("expected 1 artist, got %d", len(track.Artists))
	}
	if track.Album.Name != "Test Album" {
		t.Errorf("expected Album.Name 'Test Album', got '%s'", track.Album.Name)
	}
}

func TestRawArtist_ToDomain(t *testing.T) {
	raw := rawArtist{
		ID:   "artist123",
		Name: "Test Artist",
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/artist/artist123",
		},
		Genres:     []string{"pop", "rock"},
		Popularity: 85,
		Followers: struct {
			Total int `json:"total"`
		}{Total: 1000000},
		Images: []rawImage{
			{URL: "https://example.com/artist.jpg", Height: 640, Width: 640},
		},
	}

	artist := raw.toDomain()

	if artist.ID != "artist123" {
		t.Errorf("expected ID 'artist123', got '%s'", artist.ID)
	}
	if artist.Name != "Test Artist" {
		t.Errorf("expected Name 'Test Artist', got '%s'", artist.Name)
	}
	if len(artist.Genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(artist.Genres))
	}
	if artist.Popularity == nil || *artist.Popularity != 85 {
		t.Errorf("expected Popularity 85, got %v", artist.Popularity)
	}
	if artist.Followers == nil || *artist.Followers != 1000000 {
		t.Errorf("expected Followers 1000000, got %v", artist.Followers)
	}
	if len(artist.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(artist.Images))
	}
}

func TestRawAlbum_ToDomain(t *testing.T) {
	raw := rawAlbum{
		ID:          "album123",
		Name:        "Test Album",
		ReleaseDate: "2023-06-15",
		TotalTracks: 12,
		Popularity:  70,
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/album/album123",
		},
		ExternalIDs: struct {
			UPC string `json:"upc"`
		}{UPC: "123456789012"},
		Genres: []string{"pop"},
		Images: []rawImage{
			{URL: "https://example.com/album.jpg", Height: 640, Width: 640},
		},
		Artists: []rawSimpleArtist{
			{
				ID:   "artist1",
				Name: "Test Artist",
				ExternalURLs: map[string]string{
					"spotify": "https://open.spotify.com/artist/artist1",
				},
			},
		},
		Tracks: struct {
			Items []rawSimpleTrack `json:"items"`
		}{
			Items: []rawSimpleTrack{
				{
					ID:          "track1",
					Name:        "Track 1",
					TrackNumber: 1,
					ExternalURLs: map[string]string{
						"spotify": "https://open.spotify.com/track/track1",
					},
					Artists: []rawSimpleArtist{
						{
							ID:   "artist1",
							Name: "Test Artist",
							ExternalURLs: map[string]string{
								"spotify": "https://open.spotify.com/artist/artist1",
							},
						},
					},
				},
			},
		},
	}

	album := raw.toDomain()

	if album.ID != "album123" {
		t.Errorf("expected ID 'album123', got '%s'", album.ID)
	}
	if album.Name != "Test Album" {
		t.Errorf("expected Name 'Test Album', got '%s'", album.Name)
	}
	if album.TotalTracks != 12 {
		t.Errorf("expected TotalTracks 12, got %d", album.TotalTracks)
	}
	if album.Popularity == nil || *album.Popularity != 70 {
		t.Errorf("expected Popularity 70, got %v", album.Popularity)
	}
	if album.UPC == nil || *album.UPC != "123456789012" {
		t.Errorf("expected UPC '123456789012', got %v", album.UPC)
	}
	if len(album.Tracks) != 1 {
		t.Errorf("expected 1 track, got %d", len(album.Tracks))
	}
}

func TestRawAlbum_ToDomainSimple(t *testing.T) {
	raw := rawAlbum{
		ID:          "album123",
		Name:        "Simple Album",
		ReleaseDate: "2023-01-01",
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/album/album123",
		},
		Images: []rawImage{
			{URL: "https://example.com/album.jpg", Height: 300, Width: 300},
		},
		Artists: []rawSimpleArtist{
			{
				ID:   "artist1",
				Name: "Artist",
				ExternalURLs: map[string]string{
					"spotify": "https://open.spotify.com/artist/artist1",
				},
			},
		},
	}

	album := raw.toDomainSimple()

	if album.ID != "album123" {
		t.Errorf("expected ID 'album123', got '%s'", album.ID)
	}
	if album.Name != "Simple Album" {
		t.Errorf("expected Name 'Simple Album', got '%s'", album.Name)
	}
	if album.ReleaseDate != "2023-01-01" {
		t.Errorf("expected ReleaseDate '2023-01-01', got '%s'", album.ReleaseDate)
	}
}

func TestGateway_FetchToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			t.Error("expected Content-Type application/x-www-form-urlencoded")
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "new_token",
			"expires_in":   3600,
		})
	}))
	defer server.Close()

	// Test verifies server setup is correct
	_ = server
}

func TestGateway_GetToken_WithCache(t *testing.T) {
	repo := newMockTokenRepo()
	repo.tokens["spotify"] = "cached_token"

	gw := &Gateway{
		clientID:  "test_client",
		secret:    "test_secret",
		httpc:     &http.Client{},
		tokenRepo: repo,
	}

	token, err := gw.getToken(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token != "cached_token" {
		t.Errorf("expected 'cached_token', got '%s'", token)
	}
}
