package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestHealthHandler_VersionInfo(t *testing.T) {
	// テスト用にバージョン情報を設定
	originalVersion := Version
	originalBuildTime := BuildTime
	originalGitCommit := GitCommit
	defer func() {
		Version = originalVersion
		BuildTime = originalBuildTime
		GitCommit = originalGitCommit
	}()

	Version = "1.0.0"
	BuildTime = "2025-12-02T12:00:00Z"
	GitCommit = "abc1234"

	h := NewHealthHandler(EnabledServices{})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.Check(w, req)

	var resp successResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	result := resp.Result.(map[string]interface{})

	if result["version"] != "1.0.0" {
		t.Errorf("version = %v, want 1.0.0", result["version"])
	}
	if result["build_time"] != "2025-12-02T12:00:00Z" {
		t.Errorf("build_time = %v, want 2025-12-02T12:00:00Z", result["build_time"])
	}
	if result["git_commit"] != "abc1234" {
		t.Errorf("git_commit = %v, want abc1234", result["git_commit"])
	}
}

func TestHealthHandler_Uptime(t *testing.T) {
	h := NewHealthHandler(EnabledServices{})

	// 少し待ってからリクエスト
	time.Sleep(10 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.Check(w, req)

	var resp successResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	result := resp.Result.(map[string]interface{})
	uptime := result["uptime"].(string)

	// uptimeが空でないことを確認
	if uptime == "" {
		t.Error("uptime should not be empty")
	}
}

func TestHealthHandler_RuntimeInfo(t *testing.T) {
	h := NewHealthHandler(EnabledServices{})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.Check(w, req)

	var resp successResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	result := resp.Result.(map[string]interface{})
	runtime := result["runtime"].(map[string]interface{})

	// Go versionは "go" で始まる
	goVersion := runtime["go_version"].(string)
	if !strings.HasPrefix(goVersion, "go") {
		t.Errorf("go_version = %s, should start with 'go'", goVersion)
	}

	// num_goroutineは1以上
	numGoroutine := runtime["num_goroutine"].(float64)
	if numGoroutine < 1 {
		t.Errorf("num_goroutine = %v, should be >= 1", numGoroutine)
	}

	// num_cpuは1以上
	numCPU := runtime["num_cpu"].(float64)
	if numCPU < 1 {
		t.Errorf("num_cpu = %v, should be >= 1", numCPU)
	}

	// goosとgoarchは空でない
	if runtime["goos"] == "" {
		t.Error("goos should not be empty")
	}
	if runtime["goarch"] == "" {
		t.Error("goarch should not be empty")
	}
}

func TestNewHealthHandler(t *testing.T) {
	services := EnabledServices{
		Spotify:      true,
		KKBOX:        true,
		Deezer:       false,
		MusicBrainz:  false,
		LastFM:       true,
		YouTubeMusic: false,
		Redis:        true,
	}

	h := NewHealthHandler(services)

	if h == nil {
		t.Fatal("NewHealthHandler returned nil")
	}

	if h.enabledServices.Spotify != true {
		t.Error("Spotify should be enabled")
	}
	if h.enabledServices.Deezer != false {
		t.Error("Deezer should be disabled")
	}
	if h.startTime.IsZero() {
		t.Error("startTime should be set")
	}
}
