package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockRedisRepo is a simple mock for testing
type mockRedisRepo struct {
	tokens         map[string]string
	saveErr        error
	getErr         error
	isValidFunc    func(key string) bool
	invalidatedKey string
	invalidateErr  error
}

func newMockRedisRepo() *mockRedisRepo {
	return &mockRedisRepo{
		tokens: make(map[string]string),
	}
}

func (m *mockRedisRepo) SaveToken(ctx context.Context, key string, token string, ttlSeconds int) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.tokens[key] = token
	return nil
}

func (m *mockRedisRepo) GetToken(ctx context.Context, key string) (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	if token, ok := m.tokens[key]; ok {
		return token, nil
	}
	return "", errors.New("not found")
}

func (m *mockRedisRepo) IsTokenValid(ctx context.Context, key string) bool {
	if m.isValidFunc != nil {
		return m.isValidFunc(key)
	}
	_, ok := m.tokens[key]
	return ok
}

func (m *mockRedisRepo) InvalidateToken(ctx context.Context, key string) error {
	if m.invalidateErr != nil {
		return m.invalidateErr
	}
	m.invalidatedKey = key
	delete(m.tokens, key)
	return nil
}

func TestCachedTokenRepository_SaveToken(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		token      string
		ttlSeconds int
		redisErr   error
		wantErr    bool
	}{
		{
			name:       "正常系: L1とL2に保存",
			key:        "spotify",
			token:      "test-token",
			ttlSeconds: 3600,
			redisErr:   nil,
			wantErr:    false,
		},
		{
			name:       "正常系: L2エラーでもL1に保存成功",
			key:        "spotify",
			token:      "test-token",
			ttlSeconds: 3600,
			redisErr:   errors.New("redis error"),
			wantErr:    false, // L1は成功するのでエラーにならない
		},
		{
			name:       "正常系: TTLが小さい場合",
			key:        "kkbox",
			token:      "test-token",
			ttlSeconds: 30,
			redisErr:   nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := newMockRedisRepo()
			mockRedis.saveErr = tt.redisErr

			repo := NewCachedTokenRepository(mockRedis)
			err := repo.SaveToken(context.Background(), tt.key, tt.token, tt.ttlSeconds)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// L1に保存されていることを確認
			repo.mu.RLock()
			entry, ok := repo.memory[tt.key]
			repo.mu.RUnlock()

			if !ok {
				t.Errorf("SaveToken() token not saved to L1")
				return
			}

			if entry.token != tt.token {
				t.Errorf("SaveToken() L1 token = %v, want %v", entry.token, tt.token)
			}
		})
	}
}

func TestCachedTokenRepository_SaveToken_NilRedis(t *testing.T) {
	repo := NewCachedTokenRepository(nil)
	err := repo.SaveToken(context.Background(), "test", "token", 3600)

	if err != nil {
		t.Errorf("SaveToken() with nil redis should not error, got %v", err)
	}

	repo.mu.RLock()
	entry, ok := repo.memory["test"]
	repo.mu.RUnlock()

	if !ok || entry.token != "token" {
		t.Errorf("SaveToken() should save to L1 even with nil redis")
	}
}

func TestCachedTokenRepository_GetToken(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		l1Token   string
		l1Expired bool
		l2Token   string
		l2Err     error
		wantToken string
		wantErr   bool
	}{
		{
			name:      "正常系: L1から取得",
			key:       "spotify",
			l1Token:   "l1-token",
			l1Expired: false,
			wantToken: "l1-token",
			wantErr:   false,
		},
		{
			name:      "正常系: L1期限切れ、L2から取得",
			key:       "spotify",
			l1Token:   "l1-token",
			l1Expired: true,
			l2Token:   "l2-token",
			wantToken: "l2-token",
			wantErr:   false,
		},
		{
			name:      "正常系: L1なし、L2から取得",
			key:       "spotify",
			l1Token:   "",
			l2Token:   "l2-token",
			wantToken: "l2-token",
			wantErr:   false,
		},
		{
			name:      "正常系: 両方なし",
			key:       "spotify",
			l1Token:   "",
			l2Token:   "",
			l2Err:     errors.New("not found"),
			wantToken: "",
			wantErr:   false, // エラーではなく空文字を返す
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := newMockRedisRepo()
			if tt.l2Token != "" {
				mockRedis.tokens[tt.key] = tt.l2Token
			}
			if tt.l2Err != nil {
				mockRedis.getErr = tt.l2Err
			}

			repo := NewCachedTokenRepository(mockRedis)

			// L1にトークンを設定
			if tt.l1Token != "" {
				repo.mu.Lock()
				expiresAt := time.Now().Add(1 * time.Hour)
				if tt.l1Expired {
					expiresAt = time.Now().Add(-1 * time.Hour) // 過去
				}
				repo.memory[tt.key] = &tokenEntry{
					token:     tt.l1Token,
					expiresAt: expiresAt,
				}
				repo.mu.Unlock()
			}

			got, err := repo.GetToken(context.Background(), tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantToken {
				t.Errorf("GetToken() = %v, want %v", got, tt.wantToken)
			}
		})
	}
}

func TestCachedTokenRepository_IsTokenValid(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		l1Token   string
		l1Expired bool
		l2Valid   bool
		want      bool
	}{
		{
			name:      "L1に有効なトークンあり",
			key:       "spotify",
			l1Token:   "token",
			l1Expired: false,
			want:      true,
		},
		{
			name:      "L1期限切れ、L2有効",
			key:       "spotify",
			l1Token:   "token",
			l1Expired: true,
			l2Valid:   true,
			want:      true,
		},
		{
			name:      "L1期限切れ、L2も無効",
			key:       "spotify",
			l1Token:   "token",
			l1Expired: true,
			l2Valid:   false,
			want:      false,
		},
		{
			name:    "両方なし",
			key:     "spotify",
			l1Token: "",
			l2Valid: false,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := newMockRedisRepo()
			mockRedis.isValidFunc = func(key string) bool {
				return tt.l2Valid
			}

			repo := NewCachedTokenRepository(mockRedis)

			if tt.l1Token != "" {
				repo.mu.Lock()
				expiresAt := time.Now().Add(1 * time.Hour)
				if tt.l1Expired {
					expiresAt = time.Now().Add(-1 * time.Hour)
				}
				repo.memory[tt.key] = &tokenEntry{
					token:     tt.l1Token,
					expiresAt: expiresAt,
				}
				repo.mu.Unlock()
			}

			got := repo.IsTokenValid(context.Background(), tt.key)

			if got != tt.want {
				t.Errorf("IsTokenValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCachedTokenRepository_PromoteToL1(t *testing.T) {
	mockRedis := newMockRedisRepo()
	mockRedis.tokens["spotify"] = "redis-token"

	repo := NewCachedTokenRepository(mockRedis)

	// L1は空
	repo.mu.RLock()
	_, ok := repo.memory["spotify"]
	repo.mu.RUnlock()
	if ok {
		t.Fatal("L1 should be empty initially")
	}

	// GetTokenでL2から取得→L1に昇格
	token, _ := repo.GetToken(context.Background(), "spotify")
	if token != "redis-token" {
		t.Errorf("GetToken() = %v, want redis-token", token)
	}

	// L1に昇格されていることを確認
	repo.mu.RLock()
	entry, ok := repo.memory["spotify"]
	repo.mu.RUnlock()

	if !ok {
		t.Error("Token should be promoted to L1")
	}
	if entry.token != "redis-token" {
		t.Errorf("L1 token = %v, want redis-token", entry.token)
	}
}
