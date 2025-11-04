package api

import "github.com/t1nyb0x/tracktaste/internal/service"

type Handler struct {
	Artist *service.ArtistService
	Track  *service.TrackService
}

func ArtistHandler(artist *service.ArtistService) *Handler {
	return &Handler{Artist: artist}
}

func TrackHandler(track *service.TrackService) *Handler {
	return &Handler{Track: track}
}
