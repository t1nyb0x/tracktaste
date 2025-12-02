// Package v2 contains V2 business logic for TrackTaste.
// V2 uses Deezer + MusicBrainz for track features and multi-source candidate collection.
package v2

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/t1nyb0x/tracktaste/internal/domain"
	"github.com/t1nyb0x/tracktaste/internal/port/external"
	"github.com/t1nyb0x/tracktaste/internal/usecase"
	"github.com/t1nyb0x/tracktaste/internal/util/logger"
)

const (
	recommendV2Timeout       = 30 * time.Second
	maxRecommendedTracksV2   = 30
	maxCandidatesV2          = 50 // After filtering (increased for multi-source)
	kkboxCandidateLimitV2    = 30 // KKBOX candidates
	lastfmCandidateLimitV2   = 30 // Last.fm candidates
	mbArtistCandidateLimitV2 = 20 // MusicBrainz artist recordings
	ytmusicCandidateLimitV2  = 25 // YouTube Music candidates
	spotifyConcurrency       = 15 // Concurrent Spotify API calls
	deezerConcurrency        = 15 // Concurrent Deezer API calls
)

// RecommendUseCase handles track recommendation using Deezer + MusicBrainz.
type RecommendUseCase struct {
	spotifyAPI     external.SpotifyAPI
	kkboxAPI       external.KKBOXAPI
	deezerAPI      external.DeezerAPI
	musicBrainzAPI external.MusicBrainzAPI
	lastfmAPI      external.LastFMAPI       // Optional: can be nil
	ytmusicAPI     external.YouTubeMusicAPI // Optional: can be nil
	calculator     *SimilarityCalculator
	genreMatcher   *usecase.GenreMatcher
}

// NewRecommendUseCase creates a new RecommendUseCase.
func NewRecommendUseCase(
	spotifyAPI external.SpotifyAPI,
	kkboxAPI external.KKBOXAPI,
	deezerAPI external.DeezerAPI,
	musicBrainzAPI external.MusicBrainzAPI,
) *RecommendUseCase {
	genreMatcher := usecase.NewGenreMatcher()
	return &RecommendUseCase{
		spotifyAPI:     spotifyAPI,
		kkboxAPI:       kkboxAPI,
		deezerAPI:      deezerAPI,
		musicBrainzAPI: musicBrainzAPI,
		calculator:     NewSimilarityCalculator(DefaultWeights(), genreMatcher),
		genreMatcher:   genreMatcher,
	}
}

// NewRecommendUseCaseWithLastFM creates a new RecommendUseCase with Last.fm support.
func NewRecommendUseCaseWithLastFM(
	spotifyAPI external.SpotifyAPI,
	kkboxAPI external.KKBOXAPI,
	deezerAPI external.DeezerAPI,
	musicBrainzAPI external.MusicBrainzAPI,
	lastfmAPI external.LastFMAPI,
) *RecommendUseCase {
	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, musicBrainzAPI)
	uc.lastfmAPI = lastfmAPI
	return uc
}

// NewRecommendUseCaseFull creates a new RecommendUseCase with all optional APIs.
func NewRecommendUseCaseFull(
	spotifyAPI external.SpotifyAPI,
	kkboxAPI external.KKBOXAPI,
	deezerAPI external.DeezerAPI,
	musicBrainzAPI external.MusicBrainzAPI,
	lastfmAPI external.LastFMAPI,
	ytmusicAPI external.YouTubeMusicAPI,
) *RecommendUseCase {
	uc := NewRecommendUseCase(spotifyAPI, kkboxAPI, deezerAPI, musicBrainzAPI)
	uc.lastfmAPI = lastfmAPI
	uc.ytmusicAPI = ytmusicAPI
	return uc
}

// GetRecommendations returns recommended tracks using Deezer + MusicBrainz features.
func (uc *RecommendUseCase) GetRecommendations(
	ctx context.Context,
	trackID string,
	mode domain.RecommendMode,
	limit int,
) (*domain.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(ctx, recommendV2Timeout)
	defer cancel()

	// Update calculator weights based on mode
	uc.calculator = NewSimilarityCalculator(WeightsForMode(mode), uc.genreMatcher)

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

	// Step 3: Collect candidate tracks from multiple sources (KKBOX + Last.fm + MusicBrainz)
	logger.Info("RecommendV2", "候補トラックを複数ソースから収集")
	candidates := uc.collectCandidatesMultiSource(ctx, track, seedFeatures)
	logger.Info("RecommendV2", fmt.Sprintf("候補トラック数: %d", len(candidates)))

	if len(candidates) == 0 {
		logger.Info("RecommendV2", "レコメンドできる曲がありませんでした")
		return &domain.RecommendResult{
			SeedTrack:    *track,
			SeedFeatures: seedFeatures,
			SeedGenres:   seedGenres,
			Items:        []domain.RecommendedTrack{},
			Mode:         mode,
		}, nil
	}

	// Step 4: Enrich candidates with Spotify + Deezer in parallel (skip MusicBrainz for speed)
	logger.Info("RecommendV2", "候補のSpotify/Deezer情報を並列取得")
	candidates, candidateFeatures := uc.enrichCandidatesParallel(ctx, candidates)

	// Step 4.5: Filter candidates by genre (remove unrelated genres)
	logger.Info("RecommendV2", fmt.Sprintf("ジャンルフィルタ前: %d件", len(candidates)))
	candidates, candidateFeatures = uc.filterByGenre(candidates, candidateFeatures, seedGenres)
	logger.Info("RecommendV2", fmt.Sprintf("ジャンルフィルタ後: %d件", len(candidates)))

	// Step 5: Calculate similarity scores and rank
	logger.Info("RecommendV2", "類似度を計算")
	recommendedTracks := uc.calculateScores(
		seedFeatures, seedArtistInfo, seedGenres,
		candidates, candidateFeatures, nil, track, // Pass seed track for same-artist/series detection
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
func (uc *RecommendUseCase) getSeedFeatures(
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
func (uc *RecommendUseCase) getArtistGenres(ctx context.Context, track *domain.Track) []string {
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
func (uc *RecommendUseCase) mergeTags(mbTags, spotifyGenres []string) []string {
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
// Deprecated: Use collectCandidatesMultiSource instead.
//
//nolint:unused // kept for potential future use or reference
func (uc *RecommendUseCase) collectCandidatesV2(
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
	if kkboxTrack == nil {
		// Track not found in KKBOX catalog (not an error)
		logger.Info("RecommendV2", "KKBOX: 曲が見つかりませんでした")
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

	// Keep up to kkboxCandidateLimitV2 (30) - will be filtered by genre later
	if len(candidates) > kkboxCandidateLimitV2 {
		candidates = candidates[:kkboxCandidateLimitV2]
	}

	return candidates
}

// collectCandidatesMultiSource collects candidate tracks from multiple sources in parallel.
// Sources: (1) KKBOX recommendations (2) Last.fm similar tracks (3) MusicBrainz artist recordings (4) YouTube Music
func (uc *RecommendUseCase) collectCandidatesMultiSource(
	ctx context.Context,
	seedTrack *domain.Track,
	seedFeatures *domain.TrackFeatures,
) []domain.Track {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Map to deduplicate by ISRC or name+artist
	seen := make(map[string]bool)
	allCandidates := make([]domain.Track, 0, 100)

	// Mark seed track as seen
	if seedTrack.ISRC != nil && *seedTrack.ISRC != "" {
		seen[*seedTrack.ISRC] = true
	}

	// Helper to add candidates with deduplication
	addCandidates := func(candidates []domain.Track, source string) {
		mu.Lock()
		defer mu.Unlock()
		added := 0
		for _, c := range candidates {
			key := ""
			if c.ISRC != nil && *c.ISRC != "" {
				key = *c.ISRC
			} else {
				// Fallback to name + artist for deduplication
				artistName := ""
				if len(c.Artists) > 0 {
					artistName = c.Artists[0].Name
				}
				key = strings.ToLower(c.Name + "|" + artistName)
			}

			if seen[key] {
				continue
			}
			seen[key] = true
			allCandidates = append(allCandidates, c)
			added++
		}
		logger.Info("RecommendV2", fmt.Sprintf("[%s] %d件追加 (重複除外後)", source, added))
	}

	// (1) KKBOX recommendations
	wg.Add(1)
	go func() {
		defer wg.Done()
		candidates := uc.collectFromKKBOX(ctx, seedTrack)
		addCandidates(candidates, "KKBOX")
	}()

	// (2) Last.fm similar tracks
	if uc.lastfmAPI != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			candidates := uc.collectFromLastFM(ctx, seedTrack)
			addCandidates(candidates, "Last.fm")
		}()
	}

	// (3) MusicBrainz artist recordings (same artist's other tracks)
	if seedFeatures != nil && seedFeatures.ArtistMBID != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			candidates := uc.collectFromMusicBrainzArtist(ctx, seedFeatures.ArtistMBID, seedTrack)
			addCandidates(candidates, "MusicBrainz")
		}()
	}

	// (4) YouTube Music similar tracks
	if uc.ytmusicAPI != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			candidates := uc.collectFromYouTubeMusic(ctx, seedTrack)
			addCandidates(candidates, "YouTubeMusic")
		}()
	}

	wg.Wait()

	logger.Info("RecommendV2", fmt.Sprintf("全ソースから合計 %d件の候補を収集", len(allCandidates)))
	return allCandidates
}

// collectFromKKBOX collects candidates from KKBOX recommendations.
func (uc *RecommendUseCase) collectFromKKBOX(ctx context.Context, seedTrack *domain.Track) []domain.Track {
	if seedTrack.ISRC == nil || *seedTrack.ISRC == "" {
		return nil
	}

	kkboxTrack, err := uc.kkboxAPI.SearchByISRC(ctx, *seedTrack.ISRC)
	if err != nil {
		logger.Warning("RecommendV2", "KKBOX ISRC検索エラー: "+err.Error())
		return nil
	}
	if kkboxTrack == nil {
		// Track not found in KKBOX catalog (not an error)
		logger.Info("RecommendV2", "KKBOX: 曲が見つかりませんでした")
		return nil
	}

	similarTracks, err := uc.kkboxAPI.GetRecommendedTracks(ctx, kkboxTrack.ID)
	if err != nil {
		logger.Warning("RecommendV2", "KKBOXレコメンド取得エラー: "+err.Error())
		return nil
	}

	candidates := make([]domain.Track, 0, len(similarTracks))
	for i, st := range similarTracks {
		if i >= kkboxCandidateLimitV2 {
			break
		}
		isrc := st.ISRC
		candidates = append(candidates, domain.Track{
			ID:   st.ID,
			Name: st.Name,
			ISRC: &isrc,
		})
	}
	return candidates
}

// collectFromLastFM collects candidates from Last.fm track.getSimilar.
func (uc *RecommendUseCase) collectFromLastFM(ctx context.Context, seedTrack *domain.Track) []domain.Track {
	if uc.lastfmAPI == nil {
		return nil
	}

	// Get artist name
	artistName := ""
	if len(seedTrack.Artists) > 0 {
		artistName = seedTrack.Artists[0].Name
	}
	if artistName == "" {
		logger.Warning("RecommendV2", "Last.fm: アーティスト名が不明")
		return nil
	}

	// Get similar tracks from Last.fm
	similarTracks, err := uc.lastfmAPI.GetSimilarTracks(ctx, artistName, seedTrack.Name, lastfmCandidateLimitV2)
	if err != nil {
		logger.Warning("RecommendV2", "Last.fm類似曲取得エラー: "+err.Error())
		return nil
	}
	if len(similarTracks) == 0 {
		logger.Info("RecommendV2", "Last.fm: 類似曲が見つかりませんでした")
		return nil
	}

	// Convert to domain.Track (will be enriched with Spotify later)
	candidates := make([]domain.Track, 0, len(similarTracks))
	for _, t := range similarTracks {
		// Create a temporary track with name/artist info
		// ISRC will be resolved via Spotify search later
		candidates = append(candidates, domain.Track{
			ID:   fmt.Sprintf("lastfm:%s:%s", t.Artist, t.Name), // Temporary ID
			Name: t.Name,
			Artists: []domain.Artist{
				{Name: t.Artist},
			},
		})
	}
	return candidates
}

// collectFromMusicBrainzArtist collects other tracks by the same artist from MusicBrainz.
func (uc *RecommendUseCase) collectFromMusicBrainzArtist(ctx context.Context, artistMBID string, seedTrack *domain.Track) []domain.Track {
	recordings, err := uc.musicBrainzAPI.GetArtistRecordings(ctx, artistMBID, mbArtistCandidateLimitV2)
	if err != nil {
		logger.Warning("RecommendV2", "MusicBrainzアーティスト曲取得エラー: "+err.Error())
		return nil
	}
	if len(recordings) == 0 {
		logger.Info("RecommendV2", "MusicBrainz: アーティストの曲が見つかりませんでした")
		return nil
	}

	candidates := make([]domain.Track, 0, len(recordings))
	for _, rec := range recordings {
		// Skip if no ISRC
		if rec.ISRC == "" {
			continue
		}
		// Skip seed track
		if seedTrack.ISRC != nil && rec.ISRC == *seedTrack.ISRC {
			continue
		}

		isrc := rec.ISRC
		candidates = append(candidates, domain.Track{
			ID:   rec.MBID, // Use MBID as temporary ID
			Name: rec.Title,
			ISRC: &isrc,
		})
	}
	return candidates
}

// collectFromYouTubeMusic collects candidates from YouTube Music similar tracks.
func (uc *RecommendUseCase) collectFromYouTubeMusic(ctx context.Context, seedTrack *domain.Track) []domain.Track {
	if uc.ytmusicAPI == nil {
		return nil
	}

	// First, search for the seed track on YouTube Music to get video ID
	artistName := ""
	if len(seedTrack.Artists) > 0 {
		artistName = seedTrack.Artists[0].Name
	}
	if artistName == "" {
		logger.Warning("RecommendV2", "YouTube Music: アーティスト名が不明")
		return nil
	}

	query := fmt.Sprintf("%s %s", artistName, seedTrack.Name)
	searchResults, err := uc.ytmusicAPI.SearchTracks(ctx, query, 1)
	if err != nil {
		logger.Warning("RecommendV2", "YouTube Music検索エラー: "+err.Error())
		return nil
	}
	if len(searchResults) == 0 {
		logger.Warning("RecommendV2", "YouTube Music: 曲が見つかりません")
		return nil
	}

	videoID := searchResults[0].VideoID
	logger.Debug("RecommendV2", fmt.Sprintf("YouTube Music: found video ID=%s for seed track", videoID))

	// Get similar tracks
	similarTracks, err := uc.ytmusicAPI.GetSimilarTracks(ctx, videoID, ytmusicCandidateLimitV2)
	if err != nil {
		logger.Warning("RecommendV2", "YouTube Music類似曲取得エラー: "+err.Error())
		return nil
	}
	if len(similarTracks) == 0 {
		logger.Info("RecommendV2", "YouTube Music: 類似曲が見つかりませんでした")
		return nil
	}

	// Convert to domain.Track (will be enriched via Spotify name search later)
	candidates := make([]domain.Track, 0, len(similarTracks))
	for _, t := range similarTracks {
		candidates = append(candidates, domain.Track{
			ID:   fmt.Sprintf("ytmusic:%s", t.VideoID), // Temporary ID
			Name: t.Title,
			Artists: []domain.Artist{
				{Name: t.Artist},
			},
		})
	}
	return candidates
}

// enrichCandidatesParallel fetches Spotify track details and Deezer features in parallel.
// This combines enrichCandidatesWithSpotify and getCandidateFeatures for better performance.
// Handles both ISRC-based and name-based (Last.fm) candidates.
func (uc *RecommendUseCase) enrichCandidatesParallel(
	ctx context.Context,
	candidates []domain.Track,
) ([]domain.Track, map[string]*domain.TrackFeatures) {
	// Separate candidates with ISRC and without ISRC (Last.fm)
	var isrcCandidates []domain.Track
	var nameCandidates []domain.Track

	for _, c := range candidates {
		if c.ISRC != nil && *c.ISRC != "" {
			isrcCandidates = append(isrcCandidates, c)
		} else if len(c.Artists) > 0 && c.Artists[0].Name != "" {
			nameCandidates = append(nameCandidates, c)
		}
	}

	// Collect ISRCs
	isrcs := make([]string, 0, len(isrcCandidates))
	for _, c := range isrcCandidates {
		isrcs = append(isrcs, *c.ISRC)
	}

	// Result containers
	enrichedTracks := make(map[string]*domain.Track)   // ISRC -> Track
	features := make(map[string]*domain.TrackFeatures) // ISRC -> Features (temporary)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 1. Fetch Spotify tracks by ISRC (parallel with semaphore)
	if len(isrcs) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem := make(chan struct{}, spotifyConcurrency)
			var innerWg sync.WaitGroup

			for _, isrc := range isrcs {
				innerWg.Add(1)
				go func(isrc string) {
					defer innerWg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					track, err := uc.spotifyAPI.SearchByISRC(ctx, isrc)
					if err != nil || track == nil {
						return
					}

					mu.Lock()
					enrichedTracks[isrc] = track
					mu.Unlock()
				}(isrc)
			}
			innerWg.Wait()
		}()
	}

	// 2. Fetch Spotify tracks by name (Last.fm candidates)
	if len(nameCandidates) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem := make(chan struct{}, spotifyConcurrency)
			var innerWg sync.WaitGroup

			for _, c := range nameCandidates {
				innerWg.Add(1)
				go func(candidate domain.Track) {
					defer innerWg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					track := uc.searchSpotifyWithFallback(ctx, candidate.Name, candidate.Artists[0].Name)
					if track == nil {
						logger.Debug("RecommendV2", fmt.Sprintf("Spotifyで見つかりませんでした: %s - %s", candidate.Artists[0].Name, candidate.Name))
						return
					}
					// Use the found track's ISRC as key
					if track.ISRC != nil && *track.ISRC != "" {
						mu.Lock()
						enrichedTracks[*track.ISRC] = track
						mu.Unlock()
					}
				}(c)
			}
			innerWg.Wait()
		}()
	}

	// 3. Fetch Deezer features (parallel batch) - only for ISRC candidates
	if len(isrcs) > 0 {
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
				features[isrc] = &domain.TrackFeatures{
					ISRC:            isrc,
					BPM:             dt.BPM,
					DurationSeconds: dt.DurationSeconds,
					Gain:            dt.Gain,
				}
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Also fetch Deezer features for Last.fm candidates that were resolved
	if len(nameCandidates) > 0 {
		var resolvedISRCs []string
		for isrc := range enrichedTracks {
			found := false
			for _, existingISRC := range isrcs {
				if isrc == existingISRC {
					found = true
					break
				}
			}
			if !found {
				resolvedISRCs = append(resolvedISRCs, isrc)
			}
		}

		if len(resolvedISRCs) > 0 {
			deezerTracks, err := uc.deezerAPI.GetTracksByISRCBatch(ctx, resolvedISRCs)
			if err == nil {
				for isrc, dt := range deezerTracks {
					features[isrc] = &domain.TrackFeatures{
						ISRC:            isrc,
						BPM:             dt.BPM,
						DurationSeconds: dt.DurationSeconds,
						Gain:            dt.Gain,
					}
				}
			}
		}
	}

	// Build final results - match Spotify tracks with Deezer features
	result := make([]domain.Track, 0, len(enrichedTracks))
	finalFeatures := make(map[string]*domain.TrackFeatures)

	for isrc, track := range enrichedTracks {
		result = append(result, *track)

		// Transfer features from ISRC-keyed to TrackID-keyed
		if f, ok := features[isrc]; ok {
			f.TrackID = track.ID
			finalFeatures[track.ID] = f

			// Use Spotify artist genres as tags (faster than MusicBrainz)
			if len(track.Artists) > 0 {
				genres, err := uc.spotifyAPI.GetArtistGenres(ctx, track.Artists[0].ID)
				if err == nil && len(genres) > 0 {
					f.Tags = genres
				}
			}
		}
	}

	return result, finalFeatures
}

// filterByGenre removes candidates with unrelated genres to improve recommendation quality.
// Keeps candidates where genre bonus >= 1.0 (exact match, same group, or related).
func (uc *RecommendUseCase) filterByGenre(
	candidates []domain.Track,
	features map[string]*domain.TrackFeatures,
	seedGenres []string,
) ([]domain.Track, map[string]*domain.TrackFeatures) {
	if len(seedGenres) == 0 {
		// No seed genres to filter by, return as-is but limit to maxCandidatesV2
		if len(candidates) > maxCandidatesV2 {
			return candidates[:maxCandidatesV2], features
		}
		return candidates, features
	}

	filtered := make([]domain.Track, 0, len(candidates))
	filteredFeatures := make(map[string]*domain.TrackFeatures)

	for _, c := range candidates {
		f := features[c.ID]
		if f == nil {
			continue
		}

		// Check genre bonus
		bonus := uc.genreMatcher.CalculateBonus(seedGenres, f.Tags)
		if bonus >= 1.0 {
			// Exact match, same group, or related - keep this candidate
			filtered = append(filtered, c)
			filteredFeatures[c.ID] = f
		} else {
			// Log filtered out candidates for debugging
			logger.Debug("RecommendV2", fmt.Sprintf("ジャンルフィルタで除外: %s (bonus=%.2f, genres=%v)", c.Name, bonus, f.Tags))
		}

		// Stop if we have enough candidates
		if len(filtered) >= maxCandidatesV2 {
			break
		}
	}

	return filtered, filteredFeatures
}

// getCandidateFeatures retrieves features for candidate tracks (legacy, kept for compatibility).
// Deprecated: Use getCandidateFeaturesParallel instead.
//
//nolint:unused // kept for potential future use or reference
func (uc *RecommendUseCase) getCandidateFeatures(
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
func (uc *RecommendUseCase) calculateScores(
	seedFeatures *domain.TrackFeatures,
	seedArtistInfo *domain.ArtistInfo,
	seedGenres []string,
	candidates []domain.Track,
	candidateFeatures map[string]*domain.TrackFeatures,
	candidateArtistInfos map[string]*domain.ArtistInfo,
	seedTrack *domain.Track,
) []domain.RecommendedTrack {
	recommendedTracks := make([]domain.RecommendedTrack, 0, len(candidates))

	// Extract seed artist IDs for same-artist detection
	seedArtistIDs := make(map[string]bool)
	seedArtistNames := make(map[string]bool)
	if seedTrack != nil {
		for _, a := range seedTrack.Artists {
			seedArtistIDs[a.ID] = true
			seedArtistNames[strings.ToLower(a.Name)] = true
		}
	}

	for _, candidate := range candidates {
		candidateFeature := candidateFeatures[candidate.ID]
		candidateArtist := candidateArtistInfos[candidate.ID]

		// Calculate similarity with bonuses
		baseSim, genreBonus, artistBonus, _ := uc.calculator.CalculateWithBonus(
			seedFeatures, candidateFeature,
			seedArtistInfo, candidateArtist,
		)

		// Get match reasons
		matchReasons := uc.calculator.MatchReasons(seedFeatures, candidateFeature)

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

		// Same artist bonus (strong boost)
		sameArtistBonus := 1.0
		for _, a := range candidate.Artists {
			if seedArtistIDs[a.ID] || seedArtistNames[strings.ToLower(a.Name)] {
				sameArtistBonus = 2.5 // Strong bonus for same artist
				matchReasons = append(matchReasons, "same_artist")
				break
			}
		}

		// Series/franchise bonus (detect related works)
		seriesBonus := 1.0
		if seedTrack != nil {
			seriesBonus, matchReasons = uc.detectSeriesMatch(seedTrack.Name, candidate.Name, matchReasons)
		}

		// Apply all bonuses
		totalBonus := genreBonus * artistBonus * sameArtistBonus * seriesBonus
		finalScore := baseSim * totalBonus

		recommendedTracks = append(recommendedTracks, domain.RecommendedTrack{
			Track:           candidate,
			SimilarityScore: baseSim,
			GenreBonus:      totalBonus,
			FinalScore:      finalScore,
			MatchReasons:    matchReasons,
			Features:        candidateFeature,
		})
	}

	return recommendedTracks
}

// detectSeriesMatch detects if two tracks belong to the same series/franchise.
func (uc *RecommendUseCase) detectSeriesMatch(seedName, candidateName string, reasons []string) (float64, []string) {
	seedLower := strings.ToLower(seedName)
	candidateLower := strings.ToLower(candidateName)

	// Check for common anime/game franchise patterns
	franchisePatterns := []struct {
		keywords []string
		name     string
	}{
		{[]string{"ラブライブ", "love live", "lovelive"}, "Love Live"},
		{[]string{"アイマス", "idolm@ster", "アイドルマスター", "cinderella", "シンデレラ", "million live", "ミリオン", "shiny colors", "シャニマス"}, "THE IDOLM@STER"},
		{[]string{"バンドリ", "bang dream", "bandori", "poppin'party", "roselia", "raise a suilen", "morfonica"}, "BanG Dream"},
		{[]string{"プロセカ", "project sekai", "プロジェクトセカイ", "初音ミク"}, "Project Sekai"},
		{[]string{"ウマ娘", "uma musume", "うまむすめ"}, "Uma Musume"},
		{[]string{"hololive", "ホロライブ"}, "hololive"},
		{[]string{"にじさんじ", "nijisanji"}, "Nijisanji"},
		{[]string{"vtuber", "ブイチューバー"}, "VTuber"},
		{[]string{"fate", "fgo", "フェイト"}, "Fate"},
		{[]string{"touhou", "東方", "幻想郷"}, "Touhou"},
		{[]string{"vocaloid", "ボカロ", "初音ミク", "鏡音リン", "鏡音レン", "巡音ルカ", "gumi", "ia"}, "Vocaloid"},
		{[]string{"けいおん", "k-on"}, "K-ON!"},
		{[]string{"ガンダム", "gundam"}, "Gundam"},
		{[]string{"マクロス", "macross"}, "Macross"},
		{[]string{"リゼロ", "re:zero", "re：zero"}, "Re:Zero"},
		{[]string{"鬼滅", "demon slayer", "kimetsu"}, "Demon Slayer"},
		{[]string{"呪術廻戦", "jujutsu kaisen"}, "Jujutsu Kaisen"},
		{[]string{"進撃の巨人", "attack on titan", "shingeki"}, "Attack on Titan"},
		{[]string{"ワンピース", "one piece"}, "One Piece"},
		{[]string{"ナルト", "naruto", "boruto"}, "Naruto"},
		{[]string{"ブリーチ", "bleach"}, "Bleach"},
		{[]string{"ドラゴンボール", "dragon ball"}, "Dragon Ball"},
		{[]string{"エヴァンゲリオン", "evangelion", "エヴァ"}, "Evangelion"},
		{[]string{"ソードアート", "sword art online", "sao"}, "SAO"},
		{[]string{"チェンソーマン", "chainsaw man"}, "Chainsaw Man"},
		{[]string{"スパイファミリー", "spy x family", "spy family"}, "SPY×FAMILY"},
		{[]string{"ぼっち・ざ・ろっく", "bocchi the rock", "bocchi"}, "Bocchi the Rock!"},
		{[]string{"推しの子", "oshi no ko"}, "Oshi no Ko"},
	}

	for _, fp := range franchisePatterns {
		seedMatch := false
		candidateMatch := false

		for _, kw := range fp.keywords {
			if strings.Contains(seedLower, kw) {
				seedMatch = true
			}
			if strings.Contains(candidateLower, kw) {
				candidateMatch = true
			}
		}

		if seedMatch && candidateMatch {
			return 2.0, append(reasons, "same_series:"+fp.name)
		}
	}

	return 1.0, reasons
}

// searchSpotifyWithFallback searches Spotify for a track with multiple fallback strategies.
// It tries progressively simpler queries if exact search fails.
func (uc *RecommendUseCase) searchSpotifyWithFallback(ctx context.Context, trackName, artistName string) *domain.Track {
	// Strategy 1: Exact search with track: and artist: filters
	sanitizedTrack := sanitizeSearchQuery(trackName)
	sanitizedArtist := sanitizeSearchQuery(artistName)

	query := fmt.Sprintf("track:%s artist:%s", sanitizedTrack, sanitizedArtist)
	tracks, err := uc.spotifyAPI.SearchTracks(ctx, query)
	if err == nil && len(tracks) > 0 {
		return &tracks[0]
	}

	// Strategy 2: Simplified track name (remove parentheses, brackets content)
	simplifiedTrack := simplifyTrackName(trackName)
	if simplifiedTrack != sanitizedTrack {
		query = fmt.Sprintf("track:%s artist:%s", simplifiedTrack, sanitizedArtist)
		tracks, err = uc.spotifyAPI.SearchTracks(ctx, query)
		if err == nil && len(tracks) > 0 {
			return &tracks[0]
		}
	}

	// Strategy 3: Free text search (artist + track name)
	query = fmt.Sprintf("%s %s", sanitizedArtist, sanitizedTrack)
	tracks, err = uc.spotifyAPI.SearchTracks(ctx, query)
	if err == nil && len(tracks) > 0 {
		// Verify the result matches the artist (fuzzy match)
		for _, t := range tracks {
			if len(t.Artists) > 0 && fuzzyMatchArtist(t.Artists[0].Name, artistName) {
				return &t
			}
		}
		// Return first result if no exact artist match
		return &tracks[0]
	}

	// Strategy 4: Simplified free text search
	if simplifiedTrack != sanitizedTrack {
		query = fmt.Sprintf("%s %s", sanitizedArtist, simplifiedTrack)
		tracks, err = uc.spotifyAPI.SearchTracks(ctx, query)
		if err == nil && len(tracks) > 0 {
			return &tracks[0]
		}
	}

	return nil
}

// sanitizeSearchQuery removes special characters that may cause search issues.
var specialCharsRegex = regexp.MustCompile(`[～〜「」『』【】（）()[\]<>《》、。・"'：:；;！!？?＆&＃#＄$％%＠@＊*＋+＝=｜|＼\\／/]`)

func sanitizeSearchQuery(s string) string {
	// Remove special characters
	s = specialCharsRegex.ReplaceAllString(s, " ")
	// Collapse multiple spaces
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	// Trim
	return strings.TrimSpace(s)
}

// simplifyTrackName removes common suffixes like (feat. ...), [Remix], etc.
var subtitlePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\s*[\(（【\[].+$`),               // Remove everything after opening bracket
	regexp.MustCompile(`\s*[-－ー]\s*.+$`),               // Remove everything after dash (common in Japanese titles)
	regexp.MustCompile(`(?i)\s*(feat\.?|ft\.?).+$`),    // Remove feat. and everything after
	regexp.MustCompile(`(?i)\s*(remix|ver\.|version)`), // Remove remix/version indicators
}

func simplifyTrackName(name string) string {
	result := name
	for _, pattern := range subtitlePatterns {
		simplified := pattern.ReplaceAllString(result, "")
		if simplified != "" && len(simplified) >= 3 {
			result = simplified
		}
	}
	return sanitizeSearchQuery(result)
}

// fuzzyMatchArtist checks if two artist names are similar.
func fuzzyMatchArtist(a, b string) bool {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	if a == b {
		return true
	}

	// Check if one contains the other
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return true
	}

	// Remove common suffixes/prefixes for comparison
	normalize := func(s string) string {
		s = strings.ReplaceAll(s, "the ", "")
		s = strings.ReplaceAll(s, " the", "")
		s = strings.ReplaceAll(s, "&", "and")
		s = strings.ReplaceAll(s, "＆", "and")
		return strings.TrimSpace(s)
	}

	return normalize(a) == normalize(b)
}
