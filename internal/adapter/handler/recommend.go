package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

// RecommendUseCase defines the interface for recommendation use cases.
// Both RecommendUseCase and RecommendUseCaseV2 implement this interface.
type RecommendUseCase interface {
	GetRecommendations(ctx context.Context, trackID string, mode domain.RecommendMode, limit int) (*domain.RecommendResult, error)
}

// RecommendHandler handles recommendation requests.
type RecommendHandler struct {
	recommendUC RecommendUseCase
}

// NewRecommendHandler creates a new RecommendHandler.
func NewRecommendHandler(recommendUC RecommendUseCase) *RecommendHandler {
	return &RecommendHandler{recommendUC: recommendUC}
}

// FetchRecommendations handles GET /v1/track/recommend.
func (h *RecommendHandler) FetchRecommendations(w http.ResponseWriter, r *http.Request) {
	logger.Info("Recommend", "リクエスト開始")

	rawURL := r.URL.Query().Get("url")
	trackID, err := extractSpotifyTrackID(rawURL)
	if err != nil {
		if e, ok := err.(*extractError); ok {
			logger.Warning("Recommend", e.Message)
			badRequest(w, e.Message, e.Code)
			return
		}
		badRequest(w, "パラメータが不正です", "INVALID_PARAM")
		return
	}

	// Parse mode parameter
	modeStr := r.URL.Query().Get("mode")
	mode := domain.ParseRecommendMode(modeStr)

	// Parse limit parameter
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 && l <= 30 {
			limit = l
		}
	}

	result, err := h.recommendUC.GetRecommendations(r.Context(), trackID, mode, limit)
	if err != nil {
		switch err {
		case domain.ErrISRCNotFound:
			badRequest(w, "ISRCが見つかりませんでした", "ISRC_NOT_FOUND")
		case domain.ErrTrackNotFound:
			notFound(w, "曲が見つかりませんでした", "TRACK_NOT_FOUND")
		default:
			logger.Error("Recommend", "API エラー: "+err.Error())
			serviceUnavailable(w, "APIで問題が発生しているようです", "SOMETHING_API_ERROR")
		}
		return
	}

	resp := convertRecommendResult(result)
	logger.Info("Recommend", "リクエスト完了")
	success(w, resp)
}

// recommendResponse is the API response structure.
type recommendResponse struct {
	SeedTrack seedTrackResult          `json:"seed_track"`
	Items     []recommendedTrackResult `json:"items"`
	Mode      string                   `json:"mode"`
}

type seedTrackResult struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Artists       []recommendArtistResult `json:"artists"`
	AudioFeatures *audioFeaturesResult    `json:"audio_features,omitempty"`
	Genres        []string                `json:"genres,omitempty"`
}

type recommendedTrackResult struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Artists         []recommendArtistResult `json:"artists"`
	Album           recommendAlbumResult    `json:"album"`
	URL             string                  `json:"url"`
	SimilarityScore float64                 `json:"similarity_score"`
	GenreBonus      float64                 `json:"genre_bonus"`
	FinalScore      float64                 `json:"final_score"`
	MatchReasons    []string                `json:"match_reasons"`
	AudioFeatures   *audioFeaturesResult    `json:"audio_features,omitempty"`
}

type recommendArtistResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type recommendAlbumResult struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Images      []imageResult `json:"images"`
	ReleaseDate string        `json:"release_date"`
}

type audioFeaturesResult struct {
	Tempo        float64 `json:"tempo"`
	Energy       float64 `json:"energy"`
	Danceability float64 `json:"danceability"`
	Valence      float64 `json:"valence"`
	Acousticness float64 `json:"acousticness"`
	// New fields from TrackFeatures (Deezer + MusicBrainz)
	BPM             float64  `json:"bpm,omitempty"`
	DurationSeconds int      `json:"duration_seconds,omitempty"`
	Gain            float64  `json:"gain,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

func convertRecommendResult(result *domain.RecommendResult) recommendResponse {
	// Convert seed track
	seedArtists := make([]recommendArtistResult, len(result.SeedTrack.Artists))
	for i, a := range result.SeedTrack.Artists {
		seedArtists[i] = recommendArtistResult{ID: a.ID, Name: a.Name, URL: a.URL}
	}

	var seedFeatures *audioFeaturesResult
	// Support both old AudioFeatures and new TrackFeatures
	if result.SeedAudioFeatures != nil {
		seedFeatures = &audioFeaturesResult{
			Tempo:        result.SeedAudioFeatures.Tempo,
			Energy:       result.SeedAudioFeatures.Energy,
			Danceability: result.SeedAudioFeatures.Danceability,
			Valence:      result.SeedAudioFeatures.Valence,
			Acousticness: result.SeedAudioFeatures.Acousticness,
		}
	} else if result.SeedFeatures != nil {
		seedFeatures = &audioFeaturesResult{
			BPM:             result.SeedFeatures.BPM,
			DurationSeconds: result.SeedFeatures.DurationSeconds,
			Gain:            result.SeedFeatures.Gain,
			Tags:            result.SeedFeatures.Tags,
		}
	}

	seedTrack := seedTrackResult{
		ID:            result.SeedTrack.ID,
		Name:          result.SeedTrack.Name,
		Artists:       seedArtists,
		AudioFeatures: seedFeatures,
		Genres:        result.SeedGenres,
	}

	// Convert recommended tracks
	items := make([]recommendedTrackResult, len(result.Items))
	for i, rt := range result.Items {
		artists := make([]recommendArtistResult, len(rt.Track.Artists))
		for j, a := range rt.Track.Artists {
			artists[j] = recommendArtistResult{ID: a.ID, Name: a.Name, URL: a.URL}
		}

		images := make([]imageResult, len(rt.Track.Album.Images))
		for j, img := range rt.Track.Album.Images {
			images[j] = imageResult{URL: img.URL, Height: img.Height, Width: img.Width}
		}

		album := recommendAlbumResult{
			ID:          rt.Track.Album.ID,
			Name:        rt.Track.Album.Name,
			URL:         rt.Track.Album.URL,
			Images:      images,
			ReleaseDate: rt.Track.Album.ReleaseDate,
		}

		var features *audioFeaturesResult
		// Support both old AudioFeatures and new TrackFeatures
		if rt.AudioFeatures != nil {
			features = &audioFeaturesResult{
				Tempo:        rt.AudioFeatures.Tempo,
				Energy:       rt.AudioFeatures.Energy,
				Danceability: rt.AudioFeatures.Danceability,
				Valence:      rt.AudioFeatures.Valence,
				Acousticness: rt.AudioFeatures.Acousticness,
			}
		} else if rt.Features != nil {
			features = &audioFeaturesResult{
				BPM:             rt.Features.BPM,
				DurationSeconds: rt.Features.DurationSeconds,
				Gain:            rt.Features.Gain,
				Tags:            rt.Features.Tags,
			}
		}

		items[i] = recommendedTrackResult{
			ID:              rt.Track.ID,
			Name:            rt.Track.Name,
			Artists:         artists,
			Album:           album,
			URL:             rt.Track.URL,
			SimilarityScore: rt.SimilarityScore,
			GenreBonus:      rt.GenreBonus,
			FinalScore:      rt.FinalScore,
			MatchReasons:    rt.MatchReasons,
			AudioFeatures:   features,
		}
	}

	return recommendResponse{
		SeedTrack: seedTrack,
		Items:     items,
		Mode:      string(result.Mode),
	}
}
