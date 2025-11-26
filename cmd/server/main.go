package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/t1nyb0x/tracktaste/internal/adapter/gateway/kkbox"
	redisGateway "github.com/t1nyb0x/tracktaste/internal/adapter/gateway/redis"
	"github.com/t1nyb0x/tracktaste/internal/adapter/gateway/spotify"
	"github.com/t1nyb0x/tracktaste/internal/adapter/handler"
	"github.com/t1nyb0x/tracktaste/internal/adapter/server"
	"github.com/t1nyb0x/tracktaste/internal/usecase"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

type config struct {
	httpAddr      string
	spotifyID     string
	spotifySecret string
	kkboxID       string
	kkboxSecret   string
}

// getProjectRoot はプロジェクトルートのパスを取得します。
// このファイル（cmd/server/main.go）から2階層上がプロジェクトルートです。
func getProjectRoot() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Join(filepath.Dir(currentFile), "..", "..")
}

func loadConfig() (*config, error) {
	// プロジェクトルートの .env を読み込む
	envPath := filepath.Join(getProjectRoot(), ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Error loading .env file from %s", envPath)
	}

	cfg := &config{
		httpAddr:      getEnv("HTTP_ADDR", ":8080"),
		spotifyID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifySecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		kkboxID:       os.Getenv("KKBOX_ID"),
		kkboxSecret:   os.Getenv("KKBOX_SECRET"),
	}

	if cfg.spotifyID == "" || cfg.spotifySecret == "" {
		return nil, fmt.Errorf("SPOTIFY credentials not set")
	}
	if cfg.kkboxID == "" || cfg.kkboxSecret == "" {
		return nil, fmt.Errorf("KKBOX credentials not set")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := redisGateway.Init(); err != nil {
		logger.Warning("Main", "Redis connection failed")
	} else {
		logger.Info("Main", "Redis connected")
	}

	tokenRepo := redisGateway.NewTokenRepository()
	spotifyGW := spotify.NewGateway(cfg.spotifyID, cfg.spotifySecret, tokenRepo)
	kkboxGW := kkbox.NewGateway(cfg.kkboxID, cfg.kkboxSecret, tokenRepo)

	trackUC := usecase.NewTrackUseCase(spotifyGW)
	artistUC := usecase.NewArtistUseCase(spotifyGW)
	albumUC := usecase.NewAlbumUseCase(spotifyGW)
	similarUC := usecase.NewSimilarTracksUseCase(spotifyGW, kkboxGW)

	trackH := handler.NewTrackHandler(trackUC, similarUC)
	artistH := handler.NewArtistHandler(artistUC)
	albumH := handler.NewAlbumHandler(albumUC)

	srv := server.New(
		server.Config{Addr: cfg.httpAddr},
		server.Handlers{Track: trackH, Artist: artistH, Album: albumH},
	)

	logger.Info("Main", fmt.Sprintf("Server starting on %s", cfg.httpAddr))
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case sig := <-quit:
		logger.Info("Main", fmt.Sprintf("Shutting down: %s", sig))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		logger.Info("Main", "Server stopped")
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Main", err.Error())
		}
	}
}
