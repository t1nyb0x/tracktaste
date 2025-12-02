package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler_Check(t *testing.T) {
	tests := []struct {
		name            string
		enabledServices EnabledServices
		wantStatus      int
		checkServices   map[string]string
	}{
		{
			name: "all services enabled",
			enabledServices: EnabledServices{
				Spotify:      true,
				KKBOX:        true,
				Deezer:       true,
				MusicBrainz:  true,
				LastFM:       true,
				YouTubeMusic: true,
				Redis:        true,
			},
			wantStatus: http.StatusOK,
			checkServices: map[string]string{
				"spotify":       "enabled",
				"kkbox":         "enabled",
				"deezer":        "enabled",
				"musicbrainz":   "enabled",
				"lastfm":        "enabled",
				"youtube_music": "enabled",
				"redis":         "enabled",
			},
		},
		{
			name: "minimal services",
			enabledServices: EnabledServices{
				Spotify:      true,
				KKBOX:        true,
				Deezer:       true,
				MusicBrainz:  true,
				LastFM:       false,
				YouTubeMusic: false,
				Redis:        false,
			},
			wantStatus: http.StatusOK,
			checkServices: map[string]string{
				"spotify":       "enabled",
				"kkbox":         "enabled",
				"lastfm":        "disabled",
				"youtube_music": "disabled",
				"redis":         "disabled",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHealthHandler(tt.enabledServices)

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			w := httptest.NewRecorder()

			h.Check(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Check() status = %d, want %d", w.Code, tt.wantStatus)
			}

			var resp successResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			result, ok := resp.Result.(map[string]interface{})
			if !ok {
				t.Fatal("result is not a map")
			}

			// Check status
			if result["status"] != "healthy" {
				t.Errorf("status = %v, want healthy", result["status"])
			}

			// Check version
			if result["version"] == nil || result["version"] == "" {
				t.Error("version should be set")
			}

			// Check uptime
			if result["uptime"] == nil || result["uptime"] == "" {
				t.Error("uptime should be set")
			}

			// Check runtime info
			runtime, ok := result["runtime"].(map[string]interface{})
			if !ok {
				t.Fatal("runtime is not a map")
			}
			if runtime["go_version"] == nil {
				t.Error("go_version should be set")
			}
			if runtime["num_goroutine"] == nil {
				t.Error("num_goroutine should be set")
			}

			// Check services
			services, ok := result["services"].(map[string]interface{})
			if !ok {
				t.Fatal("services is not a map")
			}
			for key, expected := range tt.checkServices {
				if services[key] != expected {
					t.Errorf("services[%s] = %v, want %s", key, services[key], expected)
				}
			}
		})
	}
}

func TestServiceStatus(t *testing.T) {
	tests := []struct {
		enabled bool
		want    string
	}{
		{true, "enabled"},
		{false, "disabled"},
	}

	for _, tt := range tests {
		got := serviceStatus(tt.enabled)
		if got != tt.want {
			t.Errorf("serviceStatus(%v) = %s, want %s", tt.enabled, got, tt.want)
		}
	}
}

func TestHealthResponse_JSON(t *testing.T) {
	h := NewHealthHandler(EnabledServices{Spotify: true})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.Check(w, req)

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}
}
