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
	maxConcurrent    = 5
	requestTimeout   = 5 * time.Second
	overallTimeout   = 30 * time.Second
	maxSimilarTracks = 30
)

type SimilarTracksUseCase struct {
	spotifyAPI external.SpotifyAPI
	kkboxAPI   external.KKBOXAPI
}

func NewSimilarTracksUseCase(spotifyAPI external.SpotifyAPI, kkboxAPI external.KKBOXAPI) *SimilarTracksUseCase {
	return &SimilarTracksUseCase{spotifyAPI: spotifyAPI, kkboxAPI: kkboxAPI}
}

func (uc *SimilarTracksUseCase) FetchSimilar(ctx context.Context, trackID string) (*domain.SimilarTracksResult, error) {
	ctx, cancel := context.WithTimeout(ctx, overallTimeout)
	defer cancel()

	logger.Info("SimilarTracks", "Spotifyからトラック情報を取得")
	track, err := uc.spotifyAPI.GetTrackByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	if track.ISRC == nil || *track.ISRC == "" {
		return nil, domain.ErrISRCNotFound
	}
	isrc := *track.ISRC

	logger.Info("SimilarTracks", "KKBOXで検索")
	kkboxTrack, err := uc.kkboxAPI.SearchByISRC(ctx, isrc)
	if err != nil {
		return nil, err
	}
	if kkboxTrack == nil {
		return nil, domain.ErrTrackNotFound
	}

	logger.Info("SimilarTracks", "KKBOXからレコメンドトラックを取得")
	recommended, err := uc.kkboxAPI.GetRecommendedTracks(ctx, kkboxTrack.ID)
	if err != nil {
		return nil, err
	}

	if len(recommended) == 0 {
		return &domain.SimilarTracksResult{Items: []domain.SimilarTrack{}}, nil
	}

	logger.Info("SimilarTracks", "KKBOXから詳細情報を取得")
	isrcList := make([]string, 0, len(recommended))
	for _, t := range recommended {
		detail, err := uc.kkboxAPI.GetTrackDetail(ctx, t.ID)
		if err != nil || detail == nil {
			continue
		}
		if detail.ISRC != "" {
			isrcList = append(isrcList, detail.ISRC)
		}
	}

	logger.Info("SimilarTracks", "Spotifyで並列検索開始")
	similarTracks := uc.searchSpotifyParallel(ctx, isrcList)
	similarTracks = removeDuplicates(similarTracks)

	sort.Slice(similarTracks, func(i, j int) bool {
		popI, popJ := 0, 0
		if similarTracks[i].Popularity != nil {
			popI = *similarTracks[i].Popularity
		}
		if similarTracks[j].Popularity != nil {
			popJ = *similarTracks[j].Popularity
		}
		return popI > popJ
	})

	if len(similarTracks) > maxSimilarTracks {
		similarTracks = similarTracks[:maxSimilarTracks]
	}

	return &domain.SimilarTracksResult{Items: similarTracks}, nil
}

func (uc *SimilarTracksUseCase) searchSpotifyParallel(ctx context.Context, isrcList []string) []domain.SimilarTrack {
	var results []domain.SimilarTrack
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for _, isrc := range isrcList {
		select {
		case <-ctx.Done():
			return results
		default:
		}

		wg.Add(1)
		go func(isrc string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			defer cancel()

			track, err := uc.spotifyAPI.SearchByISRC(reqCtx, isrc)
			if err != nil || track == nil {
				return
			}

			result := domain.SimilarTrack{
				ID:          track.ID,
				Name:        track.Name,
				ISRC:        track.ISRC,
				URL:         track.URL,
				Popularity:  track.Popularity,
				TrackNumber: track.TrackNumber,
				Album:       track.Album,
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(isrc)
	}

	wg.Wait()
	return results
}

func removeDuplicates(tracks []domain.SimilarTrack) []domain.SimilarTrack {
	seen := make(map[string]bool)
	result := make([]domain.SimilarTrack, 0, len(tracks))
	for _, track := range tracks {
		key := track.ID
		if track.ISRC != nil && *track.ISRC != "" {
			key = *track.ISRC
		}
		if !seen[key] {
			seen[key] = true
			result = append(result, track)
		}
	}
	return result
}
