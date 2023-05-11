package SuricataMatch

import (
	ac "github.com/BobuSumisu/aho-corasick"
	"math"
)

// BulidTrie This is a function to build a trie tree
func BulidTrie(rulesinfo []RuleInfo) *ac.Trie {
	Trie := ac.NewTrieBuilder()
	for _, rule := range rulesinfo {
		Trie.AddPattern(rule.ContentPattern)
	}
	return Trie.Build()
}

func removeZeroContent(data []*ac.Match) []*ac.Match {
	var result []*ac.Match
	for _, d := range data {
		if d.MatchString() != "\x00\x00\x00" {
			result = append(result, d)
		}
	}
	return result
}

func calcScore(length int) float64 {
	return 1 + math.Pow(float64(length), math.E/4)
}

// CalcTotalScore This is a function to calculate the total score of a packet
func calcTotalScore(matches []*ac.Match) float64 {
	var baseScore float64
	for _, match := range matches {
		baseScore += calcScore(len(match.Match()))
	}
	totalScore := baseScore * math.Pow(math.E/2, float64(len(matches)-1))
	if totalScore > 20 {
		totalScore = 20
	}
	return totalScore
}
