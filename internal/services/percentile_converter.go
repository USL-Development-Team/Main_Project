package services

import (
	"math"
	"sort"
	"usl-server/internal/config"
)

// RankRange represents MMR range for a rank
type RankRange struct {
	Rank   string
	MinMMR float64
	MaxMMR float64
}

// RankInfo contains cached rank information for O(1) lookup
type RankInfo struct {
	CumulativeBelow float64
	RankPercent     float64
}

// PercentileConverter handles MMR to percentile conversions
// Exact port of JavaScript PercentileConverter with performance optimizations
type PercentileConverter struct {
	config                 *config.Config
	sortedRangesByPlaylist map[string][]RankRange
	rankOrderCache         map[string]map[string]RankInfo
}

// NewPercentileConverter creates a new percentile converter
func NewPercentileConverter(cfg *config.Config) *PercentileConverter {
	return &PercentileConverter{
		config:                 cfg,
		sortedRangesByPlaylist: make(map[string][]RankRange),
		rankOrderCache:         make(map[string]map[string]RankInfo),
	}
}

// MMRToPercentile converts MMR to percentile for a given playlist
// Exact port of JavaScript mmrToPercentile() with binary search optimization
func (p *PercentileConverter) MMRToPercentile(mmr float64, playlist string) float64 {

	// Validate inputs
	if mmr < 0 {
		mmr = 0
	}

	// For now, use hardcoded rank distributions matching JavaScript config
	// TODO: Move to config file once we have the full PercentileConfig structure
	rankDistributions := p.getRankDistributions()
	mmrRanges := p.getMMRRanges()

	if _, exists := rankDistributions[playlist]; !exists {
		playlist = "doubles" // Default fallback
	}

	distribution := rankDistributions[playlist]
	ranges := mmrRanges[playlist]

	// Build sorted ranges cache for binary search
	if _, exists := p.sortedRangesByPlaylist[playlist]; !exists {
		var sortedRanges []RankRange
		for rank, mmrRange := range ranges {
			sortedRanges = append(sortedRanges, RankRange{
				Rank:   rank,
				MinMMR: mmrRange[0],
				MaxMMR: mmrRange[1],
			})
		}
		sort.Slice(sortedRanges, func(i, j int) bool {
			return sortedRanges[i].MinMMR < sortedRanges[j].MinMMR
		})
		p.sortedRangesByPlaylist[playlist] = sortedRanges
	}

	// Binary search to find rank
	targetRank := p.binarySearchRank(mmr, p.sortedRangesByPlaylist[playlist])

	// Handle edge cases
	if targetRank == "" {
		if mmr < p.getLowestMMR(ranges) {
			return 0.00001 // Below lowest rank
		} else {
			return 99.99999 // Above highest rank
		}
	}

	// Build rank order cache for O(1) lookup
	if _, exists := p.rankOrderCache[playlist]; !exists {
		p.rankOrderCache[playlist] = p.buildRankOrderCache(playlist, distribution)
	}

	rankInfo, exists := p.rankOrderCache[playlist][targetRank]
	if !exists {
		return 50.0 // Default to median
	}

	// Calculate cumulative percentile
	cumulativePercent := rankInfo.CumulativeBelow

	// Interpolate within the rank based on MMR position
	rankRange := ranges[targetRank]
	rankSpan := rankRange[1] - rankRange[0]
	if rankSpan > 0 {
		positionWithinRank := (mmr - rankRange[0]) / rankSpan
		positionWithinRank = math.Max(0, math.Min(1, positionWithinRank))
		cumulativePercent += positionWithinRank * rankInfo.RankPercent
	}

	return math.Max(0.00001, math.Min(99.99999, cumulativePercent))
}

// MMRToNormalizedSkill converts MMR to normalized skill (0-100 scale)
// Exact port of JavaScript mmrToNormalizedSkill()
func (p *PercentileConverter) MMRToNormalizedSkill(mmr float64, playlist string) float64 {
	percentile := p.MMRToPercentile(mmr, playlist)
	return p.PercentileToNormalizedSkill(percentile)
}

// PercentileToNormalizedSkill converts percentile to normalized skill using power expansion transform
// EXACT port of JavaScript percentileToNormalizedSkill() with piecewise function
func (p *PercentileConverter) PercentileToNormalizedSkill(percentile float64) float64 {
	// POWER EXPANSION TRANSFORM v2.0 - Exact JavaScript port
	// Validate input with high precision bounds
	percentile = math.Max(0.00001, math.Min(99.99999, percentile))

	// Power expansion transform - provides superior upper-end separation
	if percentile < 85 {
		// Compress lower tiers slightly for stability
		return percentile * 0.85
	} else {
		// Exponential expansion for elite players (85%+)
		excess := percentile - 85
		baseValue := 85 * 0.85 // 72.25
		expandedExcess := math.Pow(excess/15, 1.5) * 27.75
		normalizedSkill := baseValue + expandedExcess

		// Ensure bounds with high precision (matches JavaScript SKILL_SCALE)
		return math.Max(0.00001, math.Min(99.99999, normalizedSkill))
	}
}

// NormalizedSkillToTrueSkillMu converts normalized skill to TrueSkill μ
// EXACT port of JavaScript normalizedSkillToTrueSkillMu() with proper config values
func (p *PercentileConverter) NormalizedSkillToTrueSkillMu(normalizedSkill float64) float64 {
	// Validate input
	normalizedSkill = math.Max(0, math.Min(100, normalizedSkill))

	// JavaScript config values - EXACT match
	// MU_SCALE: { MIN: 0, MAX: 2000 }
	// SKILL_SCALE: { MIN: 0.00001, MAX: 99.99999 }
	muMin := 0.0
	muMax := 2000.0
	skillMin := 0.00001
	skillMax := 99.99999

	// Linear mapping from normalized skill to μ scale (matches JavaScript exactly)
	muRange := muMax - muMin
	skillRange := skillMax - skillMin

	mu := muMin + (normalizedSkill/skillRange)*muRange

	// Ensure bounds
	return math.Max(muMin, math.Min(muMax, mu))
}

// AggregatePlaylistSkills aggregates skills across playlists with weights
// Exact port of JavaScript aggregatePlaylistSkills()
func (p *PercentileConverter) AggregatePlaylistSkills(playlistSkills map[string]*float64, weights map[string]float64) float64 {
	var totalWeightedSkill, totalWeight float64

	for playlist, skill := range playlistSkills {
		if skill != nil && weights[playlist] > 0 {
			totalWeightedSkill += *skill * weights[playlist]
			totalWeight += weights[playlist]
		}
	}

	if totalWeight == 0 {
		return 50.0 // Default to median skill if no valid data
	}

	return totalWeightedSkill / totalWeight
}

// Helper functions

func (p *PercentileConverter) binarySearchRank(mmr float64, sortedRanges []RankRange) string {
	left, right := 0, len(sortedRanges)-1

	for left <= right {
		mid := (left + right) / 2
		r := sortedRanges[mid]

		if mmr >= r.MinMMR && mmr <= r.MaxMMR {
			return r.Rank
		} else if mmr < r.MinMMR {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return "" // Not found
}

func (p *PercentileConverter) getLowestMMR(ranges map[string][2]float64) float64 {
	lowest := math.Inf(1)
	for _, mmrRange := range ranges {
		if mmrRange[0] < lowest {
			lowest = mmrRange[0]
		}
	}
	return lowest
}

func (p *PercentileConverter) buildRankOrderCache(playlist string, distribution map[string]float64) map[string]RankInfo {
	cache := make(map[string]RankInfo)

	// EXACT JavaScript rank order matching PercentileConfig.js
	var cumulative float64

	// Rank order matching JavaScript _buildRankOrderCache()
	rankOrder := []string{
		"Bronze 1", "Bronze 2", "Bronze 3",
		"Silver 1", "Silver 2", "Silver 3",
		"Gold 1", "Gold 2", "Gold 3",
		"Platinum 1", "Platinum 2", "Platinum 3",
		"Diamond 1", "Diamond 2", "Diamond 3",
		"Champion 1", "Champion 2", "Champion 3",
		"Grand Champion 1", "Grand Champion 2", "Grand Champion 3",
		"Supersonic Legend",
	}

	for _, rank := range rankOrder {
		if percent, exists := distribution[rank]; exists {
			cache[rank] = RankInfo{
				CumulativeBelow: cumulative,
				RankPercent:     percent,
			}
			cumulative += percent
		}
	}

	return cache
}

// Hardcoded configuration data (matches JavaScript PercentileConfig)
// TODO: Move these to proper config files

func (p *PercentileConverter) getRankDistributions() map[string]map[string]float64 {
	// EXACT JavaScript rank distributions from PercentileConfig.js (Season 14, 2024)
	return map[string]map[string]float64{
		"soloDuel": {
			"Bronze 1": 0.063, "Bronze 2": 0.296, "Bronze 3": 0.952,
			"Silver 1": 2.248, "Silver 2": 4.383, "Silver 3": 7.353,
			"Gold 1": 11.090, "Gold 2": 14.354, "Gold 3": 16.356,
			"Platinum 1": 16.361, "Platinum 2": 11.923, "Platinum 3": 7.116,
			"Diamond 1": 3.828, "Diamond 2": 1.864, "Diamond 3": 0.921,
			"Champion 1": 0.473, "Champion 2": 0.217, "Champion 3": 0.103,
			"Grand Champion 1": 0.053, "Grand Champion 2": 0.024, "Grand Champion 3": 0.011,
			"Supersonic Legend": 0.013,
		},
		"doubles": {
			"Bronze 1": 0.292, "Bronze 2": 0.713, "Bronze 3": 1.485,
			"Silver 1": 2.741, "Silver 2": 4.411, "Silver 3": 6.346,
			"Gold 1": 8.427, "Gold 2": 9.790, "Gold 3": 10.237,
			"Platinum 1": 10.422, "Platinum 2": 9.093, "Platinum 3": 7.552,
			"Diamond 1": 8.364, "Diamond 2": 6.109, "Diamond 3": 4.451,
			"Champion 1": 4.663, "Champion 2": 2.397, "Champion 3": 1.272,
			"Grand Champion 1": 0.809, "Grand Champion 2": 0.293, "Grand Champion 3": 0.087,
			"Supersonic Legend": 0.045,
		},
		"standard": {
			"Bronze 1": 0.112, "Bronze 2": 0.347, "Bronze 3": 0.956,
			"Silver 1": 2.316, "Silver 2": 4.882, "Silver 3": 8.466,
			"Gold 1": 12.146, "Gold 2": 13.673, "Gold 3": 12.832,
			"Platinum 1": 11.137, "Platinum 2": 8.700, "Platinum 3": 6.701,
			"Diamond 1": 6.651, "Diamond 2": 4.291, "Diamond 3": 2.742,
			"Champion 1": 2.339, "Champion 2": 0.989, "Champion 3": 0.431,
			"Grand Champion 1": 0.205, "Grand Champion 2": 0.064, "Grand Champion 3": 0.017,
			"Supersonic Legend": 0.003,
		},
	}
}

func (p *PercentileConverter) getMMRRanges() map[string]map[string][2]float64 {
	// EXACT JavaScript MMR ranges from PercentileConfig.js
	return map[string]map[string][2]float64{
		"soloDuel": {
			"Bronze 1": {0, 152}, "Bronze 2": {153, 214}, "Bronze 3": {215, 274},
			"Silver 1": {275, 334}, "Silver 2": {335, 394}, "Silver 3": {395, 454},
			"Gold 1": {455, 514}, "Gold 2": {515, 574}, "Gold 3": {575, 634},
			"Platinum 1": {635, 694}, "Platinum 2": {695, 754}, "Platinum 3": {755, 814},
			"Diamond 1": {815, 874}, "Diamond 2": {875, 934}, "Diamond 3": {935, 994},
			"Champion 1": {995, 1054}, "Champion 2": {1055, 1114}, "Champion 3": {1115, 1174},
			"Grand Champion 1": {1175, 1234}, "Grand Champion 2": {1235, 1294}, "Grand Champion 3": {1295, 1354},
			"Supersonic Legend": {1355, 2000},
		},
		"doubles": {
			"Bronze 1": {0, 152}, "Bronze 2": {153, 214}, "Bronze 3": {215, 274},
			"Silver 1": {275, 334}, "Silver 2": {335, 394}, "Silver 3": {395, 454},
			"Gold 1": {455, 514}, "Gold 2": {515, 574}, "Gold 3": {575, 634},
			"Platinum 1": {635, 694}, "Platinum 2": {695, 754}, "Platinum 3": {755, 814},
			"Diamond 1": {815, 874}, "Diamond 2": {875, 934}, "Diamond 3": {935, 994},
			"Champion 1": {995, 1074}, "Champion 2": {1075, 1174}, "Champion 3": {1175, 1274},
			"Grand Champion 1": {1275, 1374}, "Grand Champion 2": {1375, 1474}, "Grand Champion 3": {1475, 1574},
			"Supersonic Legend": {1575, 2300},
		},
		"standard": {
			"Bronze 1": {0, 152}, "Bronze 2": {153, 214}, "Bronze 3": {215, 274},
			"Silver 1": {275, 334}, "Silver 2": {335, 394}, "Silver 3": {395, 454},
			"Gold 1": {455, 514}, "Gold 2": {515, 574}, "Gold 3": {575, 634},
			"Platinum 1": {635, 694}, "Platinum 2": {695, 754}, "Platinum 3": {755, 834},
			"Diamond 1": {835, 914}, "Diamond 2": {915, 994}, "Diamond 3": {995, 1074},
			"Champion 1": {1075, 1174}, "Champion 2": {1175, 1274}, "Champion 3": {1275, 1374},
			"Grand Champion 1": {1375, 1474}, "Grand Champion 2": {1475, 1574}, "Grand Champion 3": {1575, 1674},
			"Supersonic Legend": {1675, 2300},
		},
	}
}
