// Package usecase contains business logic for TrackTaste.
package usecase

import "strings"

// GenreGroup represents a group of related genres.
type GenreGroup string

const (
	GenreGroupOtaku GenreGroup = "otaku"
	GenreGroupJPop  GenreGroup = "jpop"
	GenreGroupRock  GenreGroup = "rock"
	GenreGroupKPop  GenreGroup = "kpop"
	GenreGroupIdol  GenreGroup = "idol"
	GenreGroupOther GenreGroup = "other"
)

// genreGroups maps genre groups to their associated Spotify genres.
var genreGroups = map[GenreGroup][]string{
	// オタク系
	GenreGroupOtaku: {
		"anime",
		"japanese vgm",
		"otacore",
		"anime rock",
		"japanese vocaloid",
		"vocaloid",
		"japanese electropop",
		"denpa",
		"touhou",
		"doujin",
		"j-pixie",
	},

	// J-POP系
	GenreGroupJPop: {
		"j-pop",
		"japanese pop",
		"japanese teen pop",
		"city pop",
		"shibuya-kei",
		"japanese r&b",
		"japanese soul",
	},

	// ロック系
	GenreGroupRock: {
		"j-rock",
		"japanese rock",
		"visual kei",
		"alternative rock",
		"japanese metal",
		"japanese punk",
		"japanese indie rock",
	},

	// K-POP系
	GenreGroupKPop: {
		"k-pop",
		"korean pop",
		"k-pop boy group",
		"k-pop girl group",
		"korean r&b",
	},

	// アイドル系
	GenreGroupIdol: {
		"japanese idol",
		"japanese idol pop",
		"johnnys",
		"akb-group",
	},
}

// genreToGroup maps individual genres to their group for fast lookup.
var genreToGroup map[string]GenreGroup

// relatedGroups defines which genre groups are considered related.
var relatedGroups = map[GenreGroup][]GenreGroup{
	GenreGroupOtaku: {GenreGroupJPop, GenreGroupRock},
	GenreGroupJPop:  {GenreGroupOtaku, GenreGroupRock, GenreGroupIdol},
	GenreGroupRock:  {GenreGroupJPop, GenreGroupOtaku},
	GenreGroupKPop:  {GenreGroupIdol},
	GenreGroupIdol:  {GenreGroupJPop, GenreGroupKPop},
}

func init() {
	genreToGroup = make(map[string]GenreGroup)
	for group, genres := range genreGroups {
		for _, genre := range genres {
			genreToGroup[strings.ToLower(genre)] = group
		}
	}
}

// GenreMatcher calculates genre bonus for recommendations.
type GenreMatcher struct{}

// NewGenreMatcher creates a new GenreMatcher.
func NewGenreMatcher() *GenreMatcher {
	return &GenreMatcher{}
}

// CalculateBonus calculates the genre bonus based on seed and candidate genres.
// Returns:
// - 1.5 for exact genre match
// - 1.3 for same genre group
// - 1.0 for related genre groups
// - 0.5 for unrelated genres (penalty)
func (m *GenreMatcher) CalculateBonus(seedGenres, candidateGenres []string) float64 {
	if len(seedGenres) == 0 || len(candidateGenres) == 0 {
		return 1.0 // No penalty if genres are unknown
	}

	// Check for exact match
	if hasExactMatch(seedGenres, candidateGenres) {
		return 1.5
	}

	seedGroup := getGenreGroup(seedGenres)
	candidateGroup := getGenreGroup(candidateGenres)

	// If both are unknown/other, no penalty
	if seedGroup == GenreGroupOther && candidateGroup == GenreGroupOther {
		return 1.0
	}

	// Same genre group
	if seedGroup == candidateGroup {
		return 1.3
	}

	// Related genre groups
	if isRelatedGroup(seedGroup, candidateGroup) {
		return 1.0
	}

	// Unrelated genres - penalty
	return 0.5
}

// IsGenreMatch checks if the genres match (exact or same group).
func (m *GenreMatcher) IsGenreMatch(seedGenres, candidateGenres []string) bool {
	if hasExactMatch(seedGenres, candidateGenres) {
		return true
	}

	seedGroup := getGenreGroup(seedGenres)
	candidateGroup := getGenreGroup(candidateGenres)

	return seedGroup == candidateGroup && seedGroup != GenreGroupOther
}

// hasExactMatch checks if any genre appears in both lists.
func hasExactMatch(genres1, genres2 []string) bool {
	set := make(map[string]struct{}, len(genres1))
	for _, g := range genres1 {
		set[strings.ToLower(g)] = struct{}{}
	}

	for _, g := range genres2 {
		if _, ok := set[strings.ToLower(g)]; ok {
			return true
		}
	}
	return false
}

// getGenreGroup determines the primary genre group from a list of genres.
// Prioritizes more specific groups (otaku > idol > jpop > rock > kpop > other).
func getGenreGroup(genres []string) GenreGroup {
	groupCounts := make(map[GenreGroup]int)

	for _, genre := range genres {
		g := strings.ToLower(genre)
		if group, ok := genreToGroup[g]; ok {
			groupCounts[group]++
		}
	}

	if len(groupCounts) == 0 {
		return GenreGroupOther
	}

	// Priority order: otaku > idol > jpop > rock > kpop
	priority := []GenreGroup{
		GenreGroupOtaku,
		GenreGroupIdol,
		GenreGroupJPop,
		GenreGroupRock,
		GenreGroupKPop,
	}

	for _, group := range priority {
		if groupCounts[group] > 0 {
			return group
		}
	}

	return GenreGroupOther
}

// isRelatedGroup checks if two genre groups are related.
func isRelatedGroup(group1, group2 GenreGroup) bool {
	if group1 == GenreGroupOther || group2 == GenreGroupOther {
		return false
	}

	related := relatedGroups[group1]
	for _, g := range related {
		if g == group2 {
			return true
		}
	}
	return false
}
