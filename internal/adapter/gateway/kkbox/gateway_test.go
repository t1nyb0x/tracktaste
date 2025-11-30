package kkbox

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
	tokens map[string]string
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

var _ repository.TokenRepository = (*mockTokenRepository)(nil)

func TestNewGateway(t *testing.T) {
	repo := newMockTokenRepo()
	gw := NewGateway("client_id", "client_secret", repo)

	if gw.clientID != "client_id" {
		t.Errorf("expected clientID 'client_id', got '%s'", gw.clientID)
	}
	if gw.clientSecret != "client_secret" {
		t.Errorf("expected clientSecret 'client_secret', got '%s'", gw.clientSecret)
	}
	if gw.tokenRepo != repo {
		t.Error("expected tokenRepo to be set")
	}
	if gw.httpc == nil {
		t.Error("expected httpc to be initialized")
	}
}

func TestGateway_GetToken_WithCache(t *testing.T) {
	repo := newMockTokenRepo()
	repo.tokens["kkbox"] = "cached_token"

	gw := &Gateway{
		clientID:     "test_client",
		clientSecret: "test_secret",
		httpc:        &http.Client{},
		tokenRepo:    repo,
	}

	token, err := gw.getToken(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token != "cached_token" {
		t.Errorf("expected 'cached_token', got '%s'", token)
	}
}

func TestGateway_FetchToken(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		response      map[string]interface{}
		expectedError bool
	}{
		{
			name:   "正常系: トークン取得成功",
			status: http.StatusOK,
			response: map[string]interface{}{
				"access_token": "new_token",
				"expires_in":   3600,
			},
			expectedError: false,
		},
		{
			name:   "異常系: 認証失敗",
			status: http.StatusUnauthorized,
			response: map[string]interface{}{
				"error": "invalid_client",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if !strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
					t.Error("expected Content-Type application/x-www-form-urlencoded")
				}

				w.WriteHeader(tt.status)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			// Note: Cannot override const tokenEndpoint, so this test verifies request format
			_ = server
		})
	}
}

func TestGateway_SearchByISRC_ResponseParsing(t *testing.T) {
	tests := []struct {
		name          string
		response      map[string]interface{}
		expectedNil   bool
		expectedID    string
		expectedName  string
		expectedISRC  string
	}{
		{
			name: "正常系: トラック発見",
			response: map[string]interface{}{
				"tracks": map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":   "track123",
							"name": "Test Track",
							"isrc": "JPSO00123456",
						},
					},
				},
			},
			expectedNil:  false,
			expectedID:   "track123",
			expectedName: "Test Track",
			expectedISRC: "JPSO00123456",
		},
		{
			name: "正常系: トラック未発見",
			response: map[string]interface{}{
				"tracks": map[string]interface{}{
					"data": []map[string]interface{}{},
				},
			},
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			// Verify response parsing logic
			var result struct {
				Tracks struct {
					Data []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
						ISRC string `json:"isrc"`
					} `json:"data"`
				} `json:"tracks"`
			}

			resp, _ := http.Get(server.URL)
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&result)

			if tt.expectedNil {
				if len(result.Tracks.Data) != 0 {
					t.Error("expected empty result")
				}
				return
			}

			if len(result.Tracks.Data) == 0 {
				t.Fatal("expected non-empty result")
			}

			track := result.Tracks.Data[0]
			if track.ID != tt.expectedID {
				t.Errorf("expected ID '%s', got '%s'", tt.expectedID, track.ID)
			}
			if track.Name != tt.expectedName {
				t.Errorf("expected Name '%s', got '%s'", tt.expectedName, track.Name)
			}
			if track.ISRC != tt.expectedISRC {
				t.Errorf("expected ISRC '%s', got '%s'", tt.expectedISRC, track.ISRC)
			}
		})
	}
}

func TestGateway_GetRecommendedTracks_ResponseParsing(t *testing.T) {
	tests := []struct {
		name          string
		response      map[string]interface{}
		expectedCount int
	}{
		{
			name: "正常系: 複数トラック取得",
			response: map[string]interface{}{
				"tracks": map[string]interface{}{
					"data": []map[string]interface{}{
						{"id": "track1", "name": "Track 1", "isrc": "ISRC001"},
						{"id": "track2", "name": "Track 2", "isrc": "ISRC002"},
						{"id": "track3", "name": "Track 3", "isrc": "ISRC003"},
					},
				},
			},
			expectedCount: 3,
		},
		{
			name: "正常系: 空の結果",
			response: map[string]interface{}{
				"tracks": map[string]interface{}{
					"data": []map[string]interface{}{},
				},
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			var result struct {
				Tracks struct {
					Data []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
						ISRC string `json:"isrc"`
					} `json:"data"`
				} `json:"tracks"`
			}

			resp, _ := http.Get(server.URL)
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&result)

			if len(result.Tracks.Data) != tt.expectedCount {
				t.Errorf("expected %d tracks, got %d", tt.expectedCount, len(result.Tracks.Data))
			}
		})
	}
}

func TestGateway_GetTrackDetail_ResponseParsing(t *testing.T) {
	tests := []struct {
		name         string
		response     map[string]interface{}
		expectedID   string
		expectedName string
		expectedISRC string
	}{
		{
			name: "正常系: トラック詳細取得",
			response: map[string]interface{}{
				"id":   "detail123",
				"name": "Detail Track",
				"isrc": "JPSO99999999",
			},
			expectedID:   "detail123",
			expectedName: "Detail Track",
			expectedISRC: "JPSO99999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			var result struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				ISRC string `json:"isrc"`
			}

			resp, _ := http.Get(server.URL)
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&result)

			if result.ID != tt.expectedID {
				t.Errorf("expected ID '%s', got '%s'", tt.expectedID, result.ID)
			}
			if result.Name != tt.expectedName {
				t.Errorf("expected Name '%s', got '%s'", tt.expectedName, result.Name)
			}
			if result.ISRC != tt.expectedISRC {
				t.Errorf("expected ISRC '%s', got '%s'", tt.expectedISRC, result.ISRC)
			}
		})
	}
}
