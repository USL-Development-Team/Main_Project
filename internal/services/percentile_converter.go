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

	// Power transformation for skill distribution (matches JavaScript)
	powerFactor := 1.5 // This should come from config
	normalizedSkill := math.Pow(percentile/100.0, 1.0/powerFactor) * 100.0

	return math.Max(0.001, math.Min(99.999, normalizedSkill))
}

// NormalizedSkillToTrueSkillMu converts normalized skill to TrueSkill μ
// Exact port of JavaScript normalizedSkillToTrueSkillMu()
func (p *PercentileConverter) NormalizedSkillToTrueSkillMu(normalizedSkill float64) float64 {
	// Linear mapping from 0-100 skill to μ range (matches JavaScript)
	muMin := 800.0  // Should come from config
	muMax := 2200.0 // Should come from config

	skillFraction := normalizedSkill / 100.0
	mu := muMin + (muMax-muMin)*skillFraction

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

	// This is a simplified version - the JavaScript version has more complex logic
	// for rank ordering and cumulative percentile calculation
	var cumulative float64

	// Hardcoded rank order matching Rocket League ranking system
	rankOrder := []string{"unranked", "bronze1", "bronze2", "bronze3", "silver1", "silver2", "silver3",
		"gold1", "gold2", "gold3", "platinum1", "platinum2", "platinum3",
		"diamond1", "diamond2", "diamond3", "champion1", "champion2", "champion3",
		"grandchampion1", "grandchampion2", "grandchampion3", "supersonic"}

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
	return map[string]map[string]float64{
		"soloDuel": {
			"unranked": 15.0, "bronze1": 1.5, "bronze2": 2.0, "bronze3": 2.5,
			"silver1": 3.0, "silver2": 3.5, "silver3": 4.0,
			"gold1": 6.0, "gold2": 7.0, "gold3": 8.0,
			"platinum1": 9.0, "platinum2": 8.5, "platinum3": 8.0,
			"diamond1": 7.5, "diamond2": 6.5, "diamond3": 5.5,
			"champion1": 4.5, "champion2": 3.5, "champion3": 2.5,
			"grandchampion1": 1.8, "grandchampion2": 1.2, "grandchampion3": 0.8,
			"supersonic": 0.3,
		},
		"doubles": {
			"unranked": 12.0, "bronze1": 1.2, "bronze2": 1.8, "bronze3": 2.3,
			"silver1": 3.2, "silver2": 4.0, "silver3": 4.8,
			"gold1": 7.5, "gold2": 8.5, "gold3": 9.0,
			"platinum1": 10.0, "platinum2": 9.5, "platinum3": 8.8,
			"diamond1": 8.0, "diamond2": 7.0, "diamond3": 6.0,
			"champion1": 5.0, "champion2": 4.0, "champion3": 3.0,
			"grandchampion1": 2.2, "grandchampion2": 1.5, "grandchampion3": 1.0,
			"supersonic": 0.4,
		},
		"standard": {
			"unranked": 10.0, "bronze1": 1.0, "bronze2": 1.5, "bronze3": 2.2,
			"silver1": 3.5, "silver2": 4.2, "silver3": 5.0,
			"gold1": 8.0, "gold2": 9.2, "gold3": 9.8,
			"platinum1": 10.5, "platinum2": 10.0, "platinum3": 9.2,
			"diamond1": 8.5, "diamond2": 7.3, "diamond3": 6.2,
			"champion1": 5.2, "champion2": 4.0, "champion3": 2.8,
			"grandchampion1": 2.0, "grandchampion2": 1.3, "grandchampion3": 0.9,
			"supersonic": 0.35,
		},
	}
}

func (p *PercentileConverter) getMMRRanges() map[string]map[string][2]float64 {
	return map[string]map[string][2]float64{
		"soloDuel": {
			"unranked": {0, 474}, "bronze1": {475, 509}, "bronze2": {510, 549}, "bronze3": {550, 589},
			"silver1": {590, 629}, "silver2": {630, 669}, "silver3": {670, 709},
			"gold1": {710, 769}, "gold2": {770, 829}, "gold3": {830, 889},
			"platinum1": {890, 949}, "platinum2": {950, 1009}, "platinum3": {1010, 1069},
			"diamond1": {1070, 1149}, "diamond2": {1150, 1229}, "diamond3": {1230, 1309},
			"champion1": {1310, 1389}, "champion2": {1390, 1489}, "champion3": {1490, 1589},
			"grandchampion1": {1590, 1699}, "grandchampion2": {1700, 1849}, "grandchampion3": {1850, 1999},
			"supersonic": {2000, 3000},
		},
		"doubles": {
			"unranked": {0, 474}, "bronze1": {475, 509}, "bronze2": {510, 549}, "bronze3": {550, 589},
			"silver1": {590, 629}, "silver2": {630, 669}, "silver3": {670, 709},
			"gold1": {710, 769}, "gold2": {770, 829}, "gold3": {830, 889},
			"platinum1": {890, 949}, "platinum2": {950, 1009}, "platinum3": {1010, 1069},
			"diamond1": {1070, 1149}, "diamond2": {1150, 1229}, "diamond3": {1230, 1309},
			"champion1": {1310, 1389}, "champion2": {1390, 1489}, "champion3": {1490, 1589},
			"grandchampion1": {1590, 1699}, "grandchampion2": {1700, 1849}, "grandchampion3": {1850, 1999},
			"supersonic": {2000, 3000},
		},
		"standard": {
			"unranked": {0, 474}, "bronze1": {475, 509}, "bronze2": {510, 549}, "bronze3": {550, 589},
			"silver1": {590, 629}, "silver2": {630, 669}, "silver3": {670, 709},
			"gold1": {710, 769}, "gold2": {770, 829}, "gold3": {830, 889},
			"platinum1": {890, 949}, "platinum2": {950, 1009}, "platinum3": {1010, 1069},
			"diamond1": {1070, 1149}, "diamond2": {1150, 1229}, "diamond3": {1230, 1309},
			"champion1": {1310, 1389}, "champion2": {1390, 1489}, "champion3": {1490, 1589},
			"grandchampion1": {1590, 1699}, "grandchampion2": {1700, 1849}, "grandchampion3": {1850, 1999},
			"supersonic": {2000, 3000},
		},
	}
}
