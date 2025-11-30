package config

import (
	"testing"
)

func TestConfig_Struct(t *testing.T) {
	cfg := Config{
		HTTP: HTTP{
			Addr: ":8080",
		},
		KKBOX: KKBOX{
			APIKey: "kkbox_api_key",
			Secret: "kkbox_secret",
		},
		Spotify: Spotify{
			APIKey: "spotify_api_key",
			Secret: "spotify_secret",
		},
	}

	// HTTP設定のテスト
	if cfg.HTTP.Addr != ":8080" {
		t.Errorf("expected HTTP.Addr ':8080', got '%s'", cfg.HTTP.Addr)
	}

	// KKBOX設定のテスト
	if cfg.KKBOX.APIKey != "kkbox_api_key" {
		t.Errorf("expected KKBOX.APIKey 'kkbox_api_key', got '%s'", cfg.KKBOX.APIKey)
	}
	if cfg.KKBOX.Secret != "kkbox_secret" {
		t.Errorf("expected KKBOX.Secret 'kkbox_secret', got '%s'", cfg.KKBOX.Secret)
	}

	// Spotify設定のテスト
	if cfg.Spotify.APIKey != "spotify_api_key" {
		t.Errorf("expected Spotify.APIKey 'spotify_api_key', got '%s'", cfg.Spotify.APIKey)
	}
	if cfg.Spotify.Secret != "spotify_secret" {
		t.Errorf("expected Spotify.Secret 'spotify_secret', got '%s'", cfg.Spotify.Secret)
	}
}

func TestHTTP_Struct(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		expected string
	}{
		{
			name:     "デフォルトポート",
			addr:     ":8080",
			expected: ":8080",
		},
		{
			name:     "カスタムポート",
			addr:     ":3000",
			expected: ":3000",
		},
		{
			name:     "ホスト指定あり",
			addr:     "localhost:8080",
			expected: "localhost:8080",
		},
		{
			name:     "空文字",
			addr:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http := HTTP{Addr: tt.addr}
			if http.Addr != tt.expected {
				t.Errorf("expected Addr '%s', got '%s'", tt.expected, http.Addr)
			}
		})
	}
}

func TestKKBOX_Struct(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		secret string
	}{
		{
			name:   "正常な設定",
			apiKey: "valid_api_key",
			secret: "valid_secret",
		},
		{
			name:   "空の設定",
			apiKey: "",
			secret: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kkbox := KKBOX{APIKey: tt.apiKey, Secret: tt.secret}
			if kkbox.APIKey != tt.apiKey {
				t.Errorf("expected APIKey '%s', got '%s'", tt.apiKey, kkbox.APIKey)
			}
			if kkbox.Secret != tt.secret {
				t.Errorf("expected Secret '%s', got '%s'", tt.secret, kkbox.Secret)
			}
		})
	}
}

func TestSpotify_Struct(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		secret string
	}{
		{
			name:   "正常な設定",
			apiKey: "valid_api_key",
			secret: "valid_secret",
		},
		{
			name:   "空の設定",
			apiKey: "",
			secret: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotify := Spotify{APIKey: tt.apiKey, Secret: tt.secret}
			if spotify.APIKey != tt.apiKey {
				t.Errorf("expected APIKey '%s', got '%s'", tt.apiKey, spotify.APIKey)
			}
			if spotify.Secret != tt.secret {
				t.Errorf("expected Secret '%s', got '%s'", tt.secret, spotify.Secret)
			}
		})
	}
}

func TestConfig_ZeroValue(t *testing.T) {
	var cfg Config

	// ゼロ値のテスト
	if cfg.HTTP.Addr != "" {
		t.Errorf("expected empty HTTP.Addr, got '%s'", cfg.HTTP.Addr)
	}
	if cfg.KKBOX.APIKey != "" {
		t.Errorf("expected empty KKBOX.APIKey, got '%s'", cfg.KKBOX.APIKey)
	}
	if cfg.KKBOX.Secret != "" {
		t.Errorf("expected empty KKBOX.Secret, got '%s'", cfg.KKBOX.Secret)
	}
	if cfg.Spotify.APIKey != "" {
		t.Errorf("expected empty Spotify.APIKey, got '%s'", cfg.Spotify.APIKey)
	}
	if cfg.Spotify.Secret != "" {
		t.Errorf("expected empty Spotify.Secret, got '%s'", cfg.Spotify.Secret)
	}
}
