package handler

import (
	"net/http"
	"runtime"
	"time"
)

// Version はビルド時に設定されるバージョン情報
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// HealthHandler はヘルスチェックのハンドラーです
type HealthHandler struct {
	startTime       time.Time
	enabledServices EnabledServices
}

// EnabledServices は有効化されているサービスの設定
type EnabledServices struct {
	Spotify      bool
	KKBOX        bool
	Deezer       bool
	MusicBrainz  bool
	LastFM       bool
	YouTubeMusic bool
	Redis        bool
}

// HealthResponse はヘルスチェックのレスポンス
type HealthResponse struct {
	Status    string       `json:"status"`
	Version   string       `json:"version"`
	BuildTime string       `json:"build_time,omitempty"`
	GitCommit string       `json:"git_commit,omitempty"`
	Uptime    string       `json:"uptime"`
	Runtime   RuntimeInfo  `json:"runtime"`
	Services  ServicesInfo `json:"services"`
}

// RuntimeInfo はGoランタイムの情報
type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
}

// ServicesInfo は接続サービスの状態
type ServicesInfo struct {
	Spotify      string `json:"spotify"`
	KKBOX        string `json:"kkbox"`
	Deezer       string `json:"deezer"`
	MusicBrainz  string `json:"musicbrainz"`
	LastFM       string `json:"lastfm"`
	YouTubeMusic string `json:"youtube_music"`
	Redis        string `json:"redis"`
}

// NewHealthHandler は新しいHealthHandlerを作成します
func NewHealthHandler(services EnabledServices) *HealthHandler {
	return &HealthHandler{
		startTime:       time.Now(),
		enabledServices: services,
	}
}

// Check はヘルスチェックを実行します
func (h *HealthHandler) Check(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(h.startTime).Round(time.Second)

	response := HealthResponse{
		Status:    "healthy",
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		Uptime:    uptime.String(),
		Runtime: RuntimeInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			GOOS:         runtime.GOOS,
			GOARCH:       runtime.GOARCH,
		},
		Services: ServicesInfo{
			Spotify:      serviceStatus(h.enabledServices.Spotify),
			KKBOX:        serviceStatus(h.enabledServices.KKBOX),
			Deezer:       serviceStatus(h.enabledServices.Deezer),
			MusicBrainz:  serviceStatus(h.enabledServices.MusicBrainz),
			LastFM:       serviceStatus(h.enabledServices.LastFM),
			YouTubeMusic: serviceStatus(h.enabledServices.YouTubeMusic),
			Redis:        serviceStatus(h.enabledServices.Redis),
		},
	}

	success(w, response)
}

func serviceStatus(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}
