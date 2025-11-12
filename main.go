package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/t1nyb0x/tracktaste/internal/api"
	"github.com/t1nyb0x/tracktaste/internal/config"
	"github.com/t1nyb0x/tracktaste/internal/infra/lastfm"
	"github.com/t1nyb0x/tracktaste/internal/infra/spotify"
	"github.com/t1nyb0x/tracktaste/internal/repository"
	"github.com/t1nyb0x/tracktaste/internal/service"
	"github.com/t1nyb0x/tracktaste/server"
)

func loadConfig() (config.Config, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Config {
		HTTP: config.HTTP{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		LastFM: config.LastFM{
			APIKey: os.Getenv("LASTFM_API_KEY"),
		},
		KKBOX: config.KKBOX{
			APIKey: os.Getenv("KKBOX_ID"),
			Secret: os.Getenv("KKBOX_SECRET"),
		},
		Spotify: config.Spotify{
			APIKey: os.Getenv("SPOTIFY_CLIENT_ID"),
			Secret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		},
	}


	if cfg.LastFM.APIKey == "" {
		return cfg, fmt.Errorf("LASTFM_API_KEY is not set")
	}
	if cfg.KKBOX.APIKey == "" || cfg.KKBOX.Secret == "" {
		return cfg, fmt.Errorf("KKBOX_API_KEY or KKBOX_SECRET is not set")
	}
	if cfg.Spotify.APIKey == "" || cfg.Spotify.Secret == "" {
		return cfg, fmt.Errorf("SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET is not set")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	lfm := lastfm.New(cfg.LastFM.APIKey)
	spotifyClient := spotify.New(cfg.Spotify.APIKey, cfg.Spotify.Secret)
	repo := struct{
		repository.ArtistRepo
		repository.TrackRepo
	}{ArtistRepo: lfm, TrackRepo: spotifyClient}

	artistSvc := service.NewArtistService(repo.ArtistRepo)
	trackSvc := service.NewTrackService(repo.TrackRepo)
	h := api.NewHandler(artistSvc, trackSvc)

	srv := server.New(server.Options{Addr: cfg.HTTP.Addr}, server.Deps{Handler: h})

	// startup
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	// shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case sig := <-quit:
		log.Println("受信:", sig.String(), "シャットダウン開始...")
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("サーバーのシャットダウンに失敗しました:", err)
		}

		log.Println("サーバー終了")
		_ = srv.Shutdown(ctx)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
