package api

import "github.com/t1nyb0x/tracktaste/internal/service"

type Handler struct {
	Artist *service.ArtistService
	Track  *service.TrackService
}

func NewHandler(artist *service.ArtistService, track *service.TrackService) *Handler {
	return &Handler{Artist: artist, Track: track}
}
