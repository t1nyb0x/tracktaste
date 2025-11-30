package domain

import (
	"errors"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrTrackNotFound",
			err:      ErrTrackNotFound,
			expected: "track not found",
		},
		{
			name:     "ErrArtistNotFound",
			err:      ErrArtistNotFound,
			expected: "artist not found",
		},
		{
			name:     "ErrAlbumNotFound",
			err:      ErrAlbumNotFound,
			expected: "album not found",
		},
		{
			name:     "ErrISRCNotFound",
			err:      ErrISRCNotFound,
			expected: "ISRC not found",
		},
		{
			name:     "ErrInvalidURL",
			err:      ErrInvalidURL,
			expected: "invalid URL",
		},
		{
			name:     "ErrEmptyQuery",
			err:      ErrEmptyQuery,
			expected: "empty query",
		},
		{
			name:     "ErrExternalAPIError",
			err:      ErrExternalAPIError,
			expected: "external API error",
		},
		{
			name:     "ErrTimeout",
			err:      ErrTimeout,
			expected: "operation timed out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("expected error message '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestDomainErrors_Is(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "ErrTrackNotFound一致",
			err:      ErrTrackNotFound,
			target:   ErrTrackNotFound,
			expected: true,
		},
		{
			name:     "ErrTrackNotFoundとErrArtistNotFoundは不一致",
			err:      ErrTrackNotFound,
			target:   ErrArtistNotFound,
			expected: false,
		},
		{
			name:     "ラップされたエラー",
			err:      ErrISRCNotFound,
			target:   ErrISRCNotFound,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("expected errors.Is to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExtractError(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		message       string
		expectedError string
	}{
		{
			name:          "空のパラメータエラー",
			code:          ErrCodeEmptyParam,
			message:       "URLが入力されていません",
			expectedError: "URLが入力されていません",
		},
		{
			name:          "Spotify以外のURLエラー",
			code:          ErrCodeNotSpotifyURL,
			message:       "SpotifyのURLを入力してください",
			expectedError: "SpotifyのURLを入力してください",
		},
		{
			name:          "異なるSpotifyリソースエラー",
			code:          ErrCodeDifferentSpotifyURL,
			message:       "TrackのURLを入力してください",
			expectedError: "TrackのURLを入力してください",
		},
		{
			name:          "無効なURLエラー",
			code:          ErrCodeInvalidURL,
			message:       "パラメータが不正です",
			expectedError: "パラメータが不正です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ExtractError{
				Code:    tt.code,
				Message: tt.message,
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error message '%s', got '%s'", tt.expectedError, err.Error())
			}
			if err.Code != tt.code {
				t.Errorf("expected code '%s', got '%s'", tt.code, err.Code)
			}
		})
	}
}

func TestExtractError_ImplementsError(t *testing.T) {
	var err error = &ExtractError{
		Code:    ErrCodeEmptyParam,
		Message: "test error",
	}

	// error インターフェースを実装していることを確認
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Error())
	}
}

func TestErrorCodes(t *testing.T) {
	// エラーコード定数のテスト
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "ErrCodeEmptyParam",
			code:     ErrCodeEmptyParam,
			expected: "EMPTY_PARAM",
		},
		{
			name:     "ErrCodeNotSpotifyURL",
			code:     ErrCodeNotSpotifyURL,
			expected: "NOT_SPOTIFY_URL",
		},
		{
			name:     "ErrCodeDifferentSpotifyURL",
			code:     ErrCodeDifferentSpotifyURL,
			expected: "DIFFERENT_SPOTIFY_URL",
		},
		{
			name:     "ErrCodeInvalidURL",
			code:     ErrCodeInvalidURL,
			expected: "INVALID_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("expected code '%s', got '%s'", tt.expected, tt.code)
			}
		})
	}
}

func TestExtractError_TypeAssertion(t *testing.T) {
	err := &ExtractError{
		Code:    ErrCodeEmptyParam,
		Message: "パラメータが空です",
	}

	// error型からExtractErrorへの型アサーション
	var genericErr error = err
	extractErr, ok := genericErr.(*ExtractError)
	if !ok {
		t.Fatal("expected type assertion to succeed")
	}

	if extractErr.Code != ErrCodeEmptyParam {
		t.Errorf("expected code '%s', got '%s'", ErrCodeEmptyParam, extractErr.Code)
	}
	if extractErr.Message != "パラメータが空です" {
		t.Errorf("expected message 'パラメータが空です', got '%s'", extractErr.Message)
	}
}
