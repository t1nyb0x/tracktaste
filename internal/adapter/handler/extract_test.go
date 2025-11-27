package handler

import (
	"testing"
)

func TestExtractSpotifyTrackID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantID   string
		wantErr  bool
		wantCode string
	}{
		// 正常系
		{
			name:    "標準URL",
			url:     "https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC",
			wantID:  "4uLU6hMCjMI75M1A2tKUQC",
			wantErr: false,
		},
		{
			name:    "intl-ja付きURL",
			url:     "https://open.spotify.com/intl-ja/track/4uLU6hMCjMI75M1A2tKUQC",
			wantID:  "4uLU6hMCjMI75M1A2tKUQC",
			wantErr: false,
		},
		{
			name:    "クエリパラメータ付きURL",
			url:     "https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC?si=abc123",
			wantID:  "4uLU6hMCjMI75M1A2tKUQC",
			wantErr: false,
		},
		{
			name:    "intl-ja + クエリパラメータ付きURL",
			url:     "https://open.spotify.com/intl-ja/track/4uLU6hMCjMI75M1A2tKUQC?si=abc123",
			wantID:  "4uLU6hMCjMI75M1A2tKUQC",
			wantErr: false,
		},
		// 異常系
		{
			name:     "空文字",
			url:      "",
			wantErr:  true,
			wantCode: "EMPTY_PARAM",
		},
		{
			name:     "Spotify以外のURL",
			url:      "https://music.apple.com/track/123",
			wantErr:  true,
			wantCode: "NOT_SPOTIFY_URL",
		},
		{
			name:     "artistのURL",
			url:      "https://open.spotify.com/artist/4uLU6hMCjMI75M1A2tKUQC",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:     "albumのURL",
			url:      "https://open.spotify.com/album/4uLU6hMCjMI75M1A2tKUQC",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:     "不正な形式",
			url:      "https://open.spotify.com/invalid/format",
			wantErr:  true,
			wantCode: "INVALID_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := extractSpotifyTrackID(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractSpotifyTrackID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if e, ok := err.(*extractError); ok {
					if e.Code != tt.wantCode {
						t.Errorf("extractSpotifyTrackID() error code = %v, want %v", e.Code, tt.wantCode)
					}
				} else {
					t.Errorf("extractSpotifyTrackID() error type = %T, want *extractError", err)
				}
				return
			}

			if err != nil {
				t.Errorf("extractSpotifyTrackID() unexpected error = %v", err)
				return
			}

			if gotID != tt.wantID {
				t.Errorf("extractSpotifyTrackID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestExtractSpotifyArtistID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantID   string
		wantErr  bool
		wantCode string
	}{
		// 正常系
		{
			name:    "標準URL",
			url:     "https://open.spotify.com/artist/0L8ExT028jH3ddEcZwqJJ5",
			wantID:  "0L8ExT028jH3ddEcZwqJJ5",
			wantErr: false,
		},
		{
			name:    "intl-ja付きURL",
			url:     "https://open.spotify.com/intl-ja/artist/0L8ExT028jH3ddEcZwqJJ5",
			wantID:  "0L8ExT028jH3ddEcZwqJJ5",
			wantErr: false,
		},
		// 異常系
		{
			name:     "trackのURL",
			url:      "https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:     "albumのURL",
			url:      "https://open.spotify.com/album/4uLU6hMCjMI75M1A2tKUQC",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := extractSpotifyArtistID(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractSpotifyArtistID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if e, ok := err.(*extractError); ok {
					if e.Code != tt.wantCode {
						t.Errorf("extractSpotifyArtistID() error code = %v, want %v", e.Code, tt.wantCode)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("extractSpotifyArtistID() unexpected error = %v", err)
				return
			}

			if gotID != tt.wantID {
				t.Errorf("extractSpotifyArtistID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestExtractSpotifyAlbumID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantID   string
		wantErr  bool
		wantCode string
	}{
		// 正常系
		{
			name:    "標準URL",
			url:     "https://open.spotify.com/album/0iiVne9c8LZC0iuhOBiTiL",
			wantID:  "0iiVne9c8LZC0iuhOBiTiL",
			wantErr: false,
		},
		{
			name:    "intl-ja付きURL",
			url:     "https://open.spotify.com/intl-ja/album/0iiVne9c8LZC0iuhOBiTiL",
			wantID:  "0iiVne9c8LZC0iuhOBiTiL",
			wantErr: false,
		},
		// 異常系
		{
			name:     "trackのURL",
			url:      "https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:     "artistのURL",
			url:      "https://open.spotify.com/artist/0L8ExT028jH3ddEcZwqJJ5",
			wantErr:  true,
			wantCode: "DIFFERENT_SPOTIFY_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := extractSpotifyAlbumID(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractSpotifyAlbumID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if e, ok := err.(*extractError); ok {
					if e.Code != tt.wantCode {
						t.Errorf("extractSpotifyAlbumID() error code = %v, want %v", e.Code, tt.wantCode)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("extractSpotifyAlbumID() unexpected error = %v", err)
				return
			}

			if gotID != tt.wantID {
				t.Errorf("extractSpotifyAlbumID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
