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
	// オタク系（アニソン、ボカロ、ゲーム音楽など）
	GenreGroupOtaku: {
		"anime",
		"anime score",
		"anime rock",
		"anison",
		"japanese vgm",
		"otacore",
		"japanese vocaloid",
		"vocaloid",
		"japanese electropop",
		"denpa",
		"touhou",
		"doujin",
		"j-pixie",
		"game soundtrack",
		"video game music",
		"japanese soundtrack",
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
		"japanese adult contemporary",
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
		"japanese emo",
		"japanese hardcore",
	},

	// K-POP系
	GenreGroupKPop: {
		"k-pop",
		"korean pop",
		"k-pop boy group",
		"k-pop girl group",
		"korean r&b",
		"k-indie",
	},

	// アイドル系
	GenreGroupIdol: {
		"japanese idol",
		"japanese idol pop",
		"johnnys",
		"akb-group",
		"idol",
	},
}

// genreToGroup maps individual genres to their group for fast lookup.
var genreToGroup map[string]GenreGroup

// relatedGroups defines which genre groups are considered related.
// More restrictive: only closely related genres get bonus.
var relatedGroups = map[GenreGroup][]GenreGroup{
	GenreGroupOtaku: {GenreGroupRock},          // Anime rock is common; J-POP/K-POP are NOT related
	GenreGroupJPop:  {GenreGroupIdol},          // J-POP and Idol overlap
	GenreGroupRock:  {GenreGroupOtaku},         // Rock and anime rock overlap
	GenreGroupKPop:  {GenreGroupIdol},          // K-POP and Idol overlap
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
// - 2.0 for exact genre match (strong boost)
// - 1.5 for same genre group
// - 1.0 for related genre groups (neutral)
// - 0.3 for unrelated genres (strong penalty)
func (m *GenreMatcher) CalculateBonus(seedGenres, candidateGenres []string) float64 {
	if len(seedGenres) == 0 || len(candidateGenres) == 0 {
		return 1.0 // No penalty if genres are unknown
	}

	// Check for exact match
	if hasExactMatch(seedGenres, candidateGenres) {
		return 2.0
	}

	seedGroup := getGenreGroup(seedGenres)
	candidateGroup := getGenreGroup(candidateGenres)

	// If both are unknown/other, no penalty
	if seedGroup == GenreGroupOther && candidateGroup == GenreGroupOther {
		return 1.0
	}

	// Same genre group
	if seedGroup == candidateGroup {
		return 1.5
	}

	// Related genre groups
	if isRelatedGroup(seedGroup, candidateGroup) {
		return 1.0
	}

	// Unrelated genres - strong penalty
	return 0.3
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
