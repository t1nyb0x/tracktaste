// Package usecase contains business logic for TrackTaste.
package usecase

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	recommendV2Timeout      = 30 * time.Second
	maxRecommendedTracksV2  = 30
	maxCandidatesV2         = 50
	kkboxCandidateLimitV2   = 30
)

// RecommendUseCaseV2 handles track recommendation using Deezer + MusicBrainz.
type RecommendUseCaseV2 struct {
	spotifyAPI     external.SpotifyAPI
	kkboxAPI       external.KKBOXAPI
	deezerAPI      external.DeezerAPI
	musicBrainzAPI external.MusicBrainzAPI
	calculatorV2   *SimilarityCalculatorV2
	genreMatcher   *GenreMatcher
}

// NewRecommendUseCaseV2 creates a new RecommendUseCaseV2.
func NewRecommendUseCaseV2(
	spotifyAPI external.SpotifyAPI,
	kkboxAPI external.KKBOXAPI,
	deezerAPI external.DeezerAPI,
	musicBrainzAPI external.MusicBrainzAPI,
) *RecommendUseCaseV2 {
	genreMatcher := NewGenreMatcher()
	return &RecommendUseCaseV2{
		spotifyAPI:     spotifyAPI,
		kkboxAPI:       kkboxAPI,
		deezerAPI:      deezerAPI,
		musicBrainzAPI: musicBrainzAPI,
		calculatorV2:   NewSimilarityCalculatorV2(DefaultWeightsV2(), genreMatcher),
		genreMatcher:   genreMatcher,
	}
}

// GetRecommendations returns recommended tracks using Deezer + MusicBrainz features.
func (uc *RecommendUseCaseV2) GetRecommendations(
	ctx context.Context,
	trackID string,
	mode domain.RecommendMode,
	limit int,
) (*domain.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(ctx, recommendV2Timeout)
	defer cancel()

	// Update calculator weights based on mode
	uc.calculatorV2 = NewSimilarityCalculatorV2(WeightsForModeV2(mode), uc.genreMatcher)

	if limit <= 0 || limit > maxRecommendedTracksV2 {
		limit = maxRecommendedTracksV2
	}

	// Step 1: Get seed track info from Spotify
	logger.Info("RecommendV2", "シードトラック情報を取得")
	track, err := uc.spotifyAPI.GetTrackByID(ctx, trackID)
	if err != nil {
		logger.Error("RecommendV2", "シードトラック取得エラー: "+err.Error())
		return nil, err
	}

	// Step 2: Get seed track features from Deezer + MusicBrainz (parallel)
	logger.Info("RecommendV2", "シードの特徴量を取得 (Deezer + MusicBrainz)")
	seedFeatures, seedArtistInfo := uc.getSeedFeatures(ctx, track)

	// Get Spotify genres for seed artist
	seedGenres := uc.getArtistGenres(ctx, track)

	// Merge tags from MusicBrainz and Spotify genres
	if seedFeatures != nil && len(seedGenres) > 0 {
		seedFeatures.Tags = uc.mergeTags(seedFeatures.Tags, seedGenres)
	}

	// Step 3: Collect candidate tracks from KKBOX
	logger.Info("RecommendV2", "候補トラックを収集")
	candidates := uc.collectCandidatesV2(ctx, track)
	logger.Info("RecommendV2", fmt.Sprintf("候補トラック数: %d", len(candidates)))

	if len(candidates) == 0 {
		return &domain.RecommendResult{
			SeedTrack:    *track,
			SeedFeatures: seedFeatures,
			SeedGenres:   seedGenres,
			Items:        []domain.RecommendedTrack{},
			Mode:         mode,
		}, nil
	}

	// Step 4: Get features for candidates (Deezer + MusicBrainz)
	logger.Info("RecommendV2", "候補の特徴量を取得")
	candidateFeatures, candidateArtistInfos := uc.getCandidateFeatures(ctx, candidates)

	// Step 5: Calculate similarity scores and rank
	logger.Info("RecommendV2", "類似度を計算")
	recommendedTracks := uc.calculateScores(
		seedFeatures, seedArtistInfo, seedGenres,
		candidates, candidateFeatures, candidateArtistInfos,
	)

	// Sort by final score (descending)
	sort.Slice(recommendedTracks, func(i, j int) bool {
		return recommendedTracks[i].FinalScore > recommendedTracks[j].FinalScore
	})

	// Limit results
	if len(recommendedTracks) > limit {
		recommendedTracks = recommendedTracks[:limit]
	}

	return &domain.RecommendResult{
		SeedTrack:    *track,
		SeedFeatures: seedFeatures,
		SeedGenres:   seedGenres,
		Items:        recommendedTracks,
		Mode:         mode,
	}, nil
}

// getSeedFeatures retrieves features for the seed track from Deezer and MusicBrainz.
func (uc *RecommendUseCaseV2) getSeedFeatures(
	ctx context.Context,
	track *domain.Track,
) (*domain.TrackFeatures, *domain.ArtistInfo) {
	features := &domain.TrackFeatures{
		TrackID: track.ID,
	}
	var artistInfo *domain.ArtistInfo

	if track.ISRC == nil || *track.ISRC == "" {
		logger.Warning("RecommendV2", "シードトラックにISRCがありません")
		return features, artistInfo
	}
	features.ISRC = *track.ISRC

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Get Deezer features
	wg.Add(1)
	go func() {
		defer wg.Done()
		deezerTrack, err := uc.deezerAPI.GetTrackByISRC(ctx, *track.ISRC)
		if err != nil {
			if err != domain.ErrNotFound {
				logger.Warning("RecommendV2", "Deezer取得エラー: "+err.Error())
			}
			return
		}
		mu.Lock()
		features.BPM = deezerTrack.BPM
		features.DurationSeconds = deezerTrack.DurationSeconds
		features.Gain = deezerTrack.Gain
		mu.Unlock()
	}()

	// Get MusicBrainz features
	wg.Add(1)
	go func() {
		defer wg.Done()
		recording, err := uc.musicBrainzAPI.GetRecordingByISRC(ctx, *track.ISRC)
		if err != nil {
			if err != domain.ErrNotFound {
				logger.Warning("RecommendV2", "MusicBrainz取得エラー: "+err.Error())
			}
			return
		}
		mu.Lock()
		features.ArtistMBID = recording.ArtistMBID
		// Convert MBTag to string slice
		tags := make([]string, len(recording.Tags))
		for i, tag := range recording.Tags {
			tags[i] = tag.Name
		}
		features.Tags = tags
		mu.Unlock()

		// Get artist relations if we have artist MBID
		if recording.ArtistMBID != "" {
			artist, err := uc.musicBrainzAPI.GetArtistWithRelations(ctx, recording.ArtistMBID)
			if err == nil {
				mu.Lock()
				artistInfo = &domain.ArtistInfo{
					MBID:      artist.MBID,
					Name:      artist.Name,
					Tags:      artist.Tags,
					Relations: artist.Relations,
				}
				mu.Unlock()
			}
		}
	}()

	wg.Wait()
	return features, artistInfo
}

// getArtistGenres gets Spotify genres for the track's primary artist.
func (uc *RecommendUseCaseV2) getArtistGenres(ctx context.Context, track *domain.Track) []string {
	if len(track.Artists) == 0 {
		return nil
	}

	artistID := track.Artists[0].ID
	genres, err := uc.spotifyAPI.GetArtistGenres(ctx, artistID)
	if err != nil {
		logger.Warning("RecommendV2", "Spotifyジャンル取得エラー: "+err.Error())
		return nil
	}
	return genres
}

// mergeTags merges MusicBrainz tags and Spotify genres, removing duplicates.
func (uc *RecommendUseCaseV2) mergeTags(mbTags, spotifyGenres []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(mbTags)+len(spotifyGenres))

	for _, tag := range mbTags {
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	for _, genre := range spotifyGenres {
		if !seen[genre] {
			seen[genre] = true
			result = append(result, genre)
		}
	}
	return result
}

// collectCandidatesV2 collects candidate tracks from KKBOX.
func (uc *RecommendUseCaseV2) collectCandidatesV2(
	ctx context.Context,
	seedTrack *domain.Track,
) []domain.Track {
	if seedTrack.ISRC == nil || *seedTrack.ISRC == "" {
		logger.Warning("RecommendV2", "ISRCがないため候補を収集できません")
		return nil
	}

	// Get KKBOX recommendations
	kkboxTrack, err := uc.kkboxAPI.SearchByISRC(ctx, *seedTrack.ISRC)
	if err != nil {
		logger.Warning("RecommendV2", "KKBOX ISRC検索エラー: "+err.Error())
		return nil
	}

	similarTracks, err := uc.kkboxAPI.GetRecommendedTracks(ctx, kkboxTrack.ID)
	if err != nil {
		logger.Warning("RecommendV2", "KKBOXレコメンド取得エラー: "+err.Error())
		return nil
	}

	// Convert KKBOXTrackInfo to Tracks and deduplicate
	seen := make(map[string]bool)
	candidates := make([]domain.Track, 0, len(similarTracks))

	for _, st := range similarTracks {
		// Skip seed track
		if st.ISRC != "" && seedTrack.ISRC != nil && st.ISRC == *seedTrack.ISRC {
			continue
		}

		// Deduplicate by ISRC
		key := st.ID
		if st.ISRC != "" {
			key = st.ISRC
		}
		if seen[key] {
			continue
		}
		seen[key] = true

		isrc := st.ISRC
		candidates = append(candidates, domain.Track{
			ID:   st.ID,
			Name: st.Name,
			ISRC: &isrc,
		})
	}

	if len(candidates) > maxCandidatesV2 {
		candidates = candidates[:maxCandidatesV2]
	}

	return candidates
}

// getCandidateFeatures retrieves features for candidate tracks.
func (uc *RecommendUseCaseV2) getCandidateFeatures(
	ctx context.Context,
	candidates []domain.Track,
) (map[string]*domain.TrackFeatures, map[string]*domain.ArtistInfo) {
	features := make(map[string]*domain.TrackFeatures)
	artistInfos := make(map[string]*domain.ArtistInfo)

	// Collect ISRCs
	isrcs := make([]string, 0, len(candidates))
	isrcToID := make(map[string]string)
	for _, c := range candidates {
		if c.ISRC != nil && *c.ISRC != "" {
			isrcs = append(isrcs, *c.ISRC)
			isrcToID[*c.ISRC] = c.ID
		}
	}

	if len(isrcs) == 0 {
		return features, artistInfos
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Get Deezer features (batch)
	wg.Add(1)
	go func() {
		defer wg.Done()
		deezerTracks, err := uc.deezerAPI.GetTracksByISRCBatch(ctx, isrcs)
		if err != nil {
			logger.Warning("RecommendV2", "Deezerバッチ取得エラー: "+err.Error())
			return
		}
		mu.Lock()
		for isrc, dt := range deezerTracks {
			trackID := isrcToID[isrc]
			if features[trackID] == nil {
				features[trackID] = &domain.TrackFeatures{TrackID: trackID, ISRC: isrc}
			}
			features[trackID].BPM = dt.BPM
			features[trackID].DurationSeconds = dt.DurationSeconds
			features[trackID].Gain = dt.Gain
		}
		mu.Unlock()
	}()

	// Get MusicBrainz features (batch - sequential due to rate limit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		recordings, err := uc.musicBrainzAPI.GetRecordingsByISRCBatch(ctx, isrcs)
		if err != nil {
			logger.Warning("RecommendV2", "MusicBrainzバッチ取得エラー: "+err.Error())
			return
		}
		mu.Lock()
		for isrc, rec := range recordings {
			trackID := isrcToID[isrc]
			if features[trackID] == nil {
				features[trackID] = &domain.TrackFeatures{TrackID: trackID, ISRC: isrc}
			}
			features[trackID].ArtistMBID = rec.ArtistMBID
			tags := make([]string, len(rec.Tags))
			for i, tag := range rec.Tags {
				tags[i] = tag.Name
			}
			features[trackID].Tags = tags
		}
		mu.Unlock()
	}()

	wg.Wait()
	return features, artistInfos
}

// calculateScores calculates similarity scores for all candidates.
func (uc *RecommendUseCaseV2) calculateScores(
	seedFeatures *domain.TrackFeatures,
	seedArtistInfo *domain.ArtistInfo,
	seedGenres []string,
	candidates []domain.Track,
	candidateFeatures map[string]*domain.TrackFeatures,
	candidateArtistInfos map[string]*domain.ArtistInfo,
) []domain.RecommendedTrack {
	recommendedTracks := make([]domain.RecommendedTrack, 0, len(candidates))

	for _, candidate := range candidates {
		candidateFeature := candidateFeatures[candidate.ID]
		candidateArtist := candidateArtistInfos[candidate.ID]

		// Calculate similarity with bonuses
		baseSim, genreBonus, artistBonus, finalScore := uc.calculatorV2.CalculateWithBonus(
			seedFeatures, candidateFeature,
			seedArtistInfo, candidateArtist,
		)

		// Get match reasons
		matchReasons := uc.calculatorV2.MatchReasonsV2(seedFeatures, candidateFeature)

		// Add genre match reason if applicable
		if candidateFeature != nil && genreBonus > 1.0 {
			if uc.genreMatcher.IsGenreMatch(seedGenres, candidateFeature.Tags) {
				matchReasons = append(matchReasons, "genre_match")
			}
		}

		// Add artist relation reason if applicable
		if artistBonus > 1.0 {
			matchReasons = append(matchReasons, "artist_relation")
		}

		recommendedTracks = append(recommendedTracks, domain.RecommendedTrack{
			Track:           candidate,
			SimilarityScore: baseSim,
			GenreBonus:      genreBonus * artistBonus, // Combine bonuses
			FinalScore:      finalScore,
			MatchReasons:    matchReasons,
			Features:        candidateFeature,
		})
	}

	return recommendedTracks
}
