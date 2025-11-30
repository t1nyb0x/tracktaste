// Package usecase contains business logic for TrackTaste.
package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	recommendTimeout      = 30 * time.Second
	maxRecommendedTracks  = 30
	maxCandidates         = 50
	spotifyCandidateLimit = 30
	kkboxCandidateLimit   = 20
)

// RecommendUseCase handles track recommendation logic.
type RecommendUseCase struct {
	spotifyAPI   external.SpotifyAPI
	kkboxAPI     external.KKBOXAPI
	calculator   *SimilarityCalculator
	genreMatcher *GenreMatcher
}

// NewRecommendUseCase creates a new RecommendUseCase.
func NewRecommendUseCase(spotifyAPI external.SpotifyAPI, kkboxAPI external.KKBOXAPI) *RecommendUseCase {
	return &RecommendUseCase{
		spotifyAPI:   spotifyAPI,
		kkboxAPI:     kkboxAPI,
		calculator:   NewSimilarityCalculator(DefaultWeights()),
		genreMatcher: NewGenreMatcher(),
	}
}

// GetRecommendations returns recommended tracks based on the input track.
func (uc *RecommendUseCase) GetRecommendations(
	ctx context.Context,
	trackID string,
	mode domain.RecommendMode,
	limit int,
) (*domain.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(ctx, recommendTimeout)
	defer cancel()

	// Update calculator weights based on mode
	uc.calculator = NewSimilarityCalculator(WeightsForMode(mode))

	if limit <= 0 || limit > maxRecommendedTracks {
		limit = maxRecommendedTracks
	}

	// Step 1: Get seed track info
	logger.Info("Recommend", "シードトラック情報を取得")
	track, err := uc.spotifyAPI.GetTrackByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	// Step 2: Get audio features for seed track
	logger.Info("Recommend", "Audio Features を取得")
	seedFeatures, err := uc.spotifyAPI.GetAudioFeatures(ctx, trackID)
	if err != nil {
		logger.Warning("Recommend", "Audio Features 取得失敗: "+err.Error())
		// Continue without audio features
	}

	// Step 3: Get seed track's artist genres
	var seedGenres []string
	if len(track.Artists) > 0 {
		seedGenres, err = uc.spotifyAPI.GetArtistGenres(ctx, track.Artists[0].ID)
		if err != nil {
			logger.Warning("Recommend", "アーティストジャンル取得失敗: "+err.Error())
		}
	}

	// Step 4: Collect candidate tracks (parallel)
	logger.Info("Recommend", "候補トラックを収集")
	candidates := uc.collectCandidates(ctx, track, seedFeatures)
	logger.Info("Recommend", "候補トラック数: "+string(rune('0'+len(candidates)/10))+string(rune('0'+len(candidates)%10)))

	if len(candidates) == 0 {
		return &domain.RecommendResult{
			SeedTrack:         *track,
			SeedAudioFeatures: seedFeatures,
			SeedGenres:        seedGenres,
			Items:             []domain.RecommendedTrack{},
			Mode:              mode,
		}, nil
	}

	// Step 5: Get audio features for candidates (batch)
	logger.Info("Recommend", "候補の Audio Features をバッチ取得")
	candidateIDs := make([]string, len(candidates))
	for i, c := range candidates {
		candidateIDs[i] = c.ID
	}
	candidateFeatures, err := uc.spotifyAPI.GetAudioFeaturesBatch(ctx, candidateIDs)
	if err != nil {
		logger.Warning("Recommend", "候補 Audio Features 取得失敗: "+err.Error())
	}

	// Create feature map for quick lookup
	featureMap := make(map[string]*domain.AudioFeatures, len(candidateFeatures))
	for i := range candidateFeatures {
		featureMap[candidateFeatures[i].TrackID] = &candidateFeatures[i]
	}

	// Step 6: Get artist genres for candidates (batch)
	logger.Info("Recommend", "候補のアーティストジャンルをバッチ取得")
	artistIDs := collectUniqueArtistIDs(candidates)
	artistGenres, err := uc.spotifyAPI.GetArtistGenresBatch(ctx, artistIDs)
	if err != nil {
		logger.Warning("Recommend", "アーティストジャンルバッチ取得失敗: "+err.Error())
		artistGenres = make(map[string][]string)
	}

	// Step 7: Calculate scores and rank
	logger.Info("Recommend", "スコア計算とランキング")
	recommendedTracks := make([]domain.RecommendedTrack, 0, len(candidates))

	for _, candidate := range candidates {
		candidateFeature := featureMap[candidate.ID]

		// Get candidate's artist genres
		var candidateGenres []string
		if len(candidate.Artists) > 0 {
			candidateGenres = artistGenres[candidate.Artists[0].ID]
		}

		// Calculate similarity score
		var similarity float64
		var matchReasons []string
		if seedFeatures != nil && candidateFeature != nil {
			similarity = uc.calculator.Calculate(seedFeatures, candidateFeature)
			matchReasons = uc.calculator.MatchReasons(seedFeatures, candidateFeature)
		} else {
			similarity = 0.5 // Default score when features are unavailable
		}

		// Calculate genre bonus
		genreBonus := uc.genreMatcher.CalculateBonus(seedGenres, candidateGenres)
		if uc.genreMatcher.IsGenreMatch(seedGenres, candidateGenres) {
			matchReasons = append(matchReasons, "same_genre")
		}

		// Final score
		finalScore := similarity * genreBonus

		recommendedTracks = append(recommendedTracks, domain.RecommendedTrack{
			Track:           candidate,
			SimilarityScore: similarity,
			GenreBonus:      genreBonus,
			FinalScore:      finalScore,
			MatchReasons:    matchReasons,
			AudioFeatures:   candidateFeature,
		})
	}

	// Sort by final score (descending)
	sort.Slice(recommendedTracks, func(i, j int) bool {
		return recommendedTracks[i].FinalScore > recommendedTracks[j].FinalScore
	})

	// Limit results
	if len(recommendedTracks) > limit {
		recommendedTracks = recommendedTracks[:limit]
	}

	return &domain.RecommendResult{
		SeedTrack:         *track,
		SeedAudioFeatures: seedFeatures,
		SeedGenres:        seedGenres,
		Items:             recommendedTracks,
		Mode:              mode,
	}, nil
}

// collectCandidates collects candidate tracks from Spotify and KKBOX in parallel.
func (uc *RecommendUseCase) collectCandidates(
	ctx context.Context,
	seedTrack *domain.Track,
	seedFeatures *domain.AudioFeatures,
) []domain.Track {
	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		candidates []domain.Track
	)

	// Spotify Recommendations API
	wg.Add(1)
	go func() {
		defer wg.Done()

		params := external.RecommendationParams{
			SeedTracks: []string{seedTrack.ID},
			Limit:      spotifyCandidateLimit,
		}

		// Add target parameters if audio features are available
		if seedFeatures != nil {
			tempo := seedFeatures.Tempo
			energy := seedFeatures.Energy
			valence := seedFeatures.Valence
			params.TargetTempo = &tempo
			params.TargetEnergy = &energy
			params.TargetValence = &valence
		}

		tracks, err := uc.spotifyAPI.GetRecommendations(ctx, params)
		if err != nil {
			logger.Warning("Recommend", "Spotify Recommendations 取得失敗: "+err.Error())
			return
		}

		mu.Lock()
		candidates = append(candidates, tracks...)
		mu.Unlock()
	}()

	// KKBOX Recommendations (via ISRC)
	if seedTrack.ISRC != nil && *seedTrack.ISRC != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Search for the track in KKBOX
			kkboxTrack, err := uc.kkboxAPI.SearchByISRC(ctx, *seedTrack.ISRC)
			if err != nil || kkboxTrack == nil {
				return
			}

			// Get KKBOX recommendations
			recommended, err := uc.kkboxAPI.GetRecommendedTracks(ctx, kkboxTrack.ID)
			if err != nil {
				logger.Warning("Recommend", "KKBOX Recommendations 取得失敗: "+err.Error())
				return
			}

			// Convert KKBOX recommendations to Spotify tracks
			kkboxTracks := uc.convertKKBOXToSpotify(ctx, recommended)

			mu.Lock()
			candidates = append(candidates, kkboxTracks...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Remove duplicates and filter
	candidates = removeDuplicateTracks(candidates, seedTrack.ID)

	// Limit candidates
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	return candidates
}

// convertKKBOXToSpotify converts KKBOX track info to Spotify tracks by ISRC lookup.
func (uc *RecommendUseCase) convertKKBOXToSpotify(
	ctx context.Context,
	kkboxTracks []external.KKBOXTrackInfo,
) []domain.Track {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		tracks []domain.Track
		sem    = make(chan struct{}, 5) // Limit concurrent requests
	)

	for _, kt := range kkboxTracks {
		if kt.ISRC == "" {
			continue
		}

		wg.Add(1)
		go func(isrc string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			track, err := uc.spotifyAPI.SearchByISRC(ctx, isrc)
			if err != nil || track == nil {
				return
			}

			mu.Lock()
			tracks = append(tracks, *track)
			mu.Unlock()
		}(kt.ISRC)
	}

	wg.Wait()
	return tracks
}

// removeDuplicateTracks removes duplicate tracks and the seed track from candidates.
func removeDuplicateTracks(tracks []domain.Track, seedID string) []domain.Track {
	seen := make(map[string]struct{})
	seen[seedID] = struct{}{} // Exclude seed track

	result := make([]domain.Track, 0, len(tracks))
	for _, t := range tracks {
		// Check by track ID
		if _, ok := seen[t.ID]; ok {
			continue
		}
		seen[t.ID] = struct{}{}

		// Also check by ISRC if available
		if t.ISRC != nil && *t.ISRC != "" {
			if _, ok := seen[*t.ISRC]; ok {
				continue
			}
			seen[*t.ISRC] = struct{}{}
		}

		result = append(result, t)
	}

	return result
}

// collectUniqueArtistIDs extracts unique artist IDs from tracks.
func collectUniqueArtistIDs(tracks []domain.Track) []string {
	seen := make(map[string]struct{})
	ids := make([]string, 0)

	for _, t := range tracks {
		for _, a := range t.Artists {
			if _, ok := seen[a.ID]; !ok {
				seen[a.ID] = struct{}{}
				ids = append(ids, a.ID)
			}
		}
	}

	// Limit to 50 (Spotify API limit)
	if len(ids) > 50 {
		ids = ids[:50]
	}

	return ids
}
