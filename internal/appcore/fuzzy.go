package appcore

import (
	"strings"
	"unicode"
)

// FuzzyScore calculates a fuzzy match score for a pattern against a target string.
// Returns a score (higher is better) and the indices of matched characters.
// Returns score of 0 if no match.
//
// Scoring algorithm:
// - Sequential character matches: +10 points each
// - Consecutive matches (streak): +15 bonus per character
// - Match at word boundary (after /, _, -, or space): +5 bonus
// - Match at start of string: +10 bonus
// - Case match: +2 bonus
func FuzzyScore(pattern, target string) (int, []int) {
	if pattern == "" {
		return 0, nil
	}

	patternLower := strings.ToLower(pattern)
	targetLower := strings.ToLower(target)

	patternRunes := []rune(patternLower)
	targetRunes := []rune(targetLower)
	targetRunesOriginal := []rune(target)

	// Find all matching positions
	var indices []int
	patternIdx := 0

	for targetIdx := 0; targetIdx < len(targetRunes) && patternIdx < len(patternRunes); targetIdx++ {
		if targetRunes[targetIdx] == patternRunes[patternIdx] {
			indices = append(indices, targetIdx)
			patternIdx++
		}
	}

	// No match if we didn't match all pattern characters
	if patternIdx < len(patternRunes) {
		return 0, nil
	}

	// Calculate score
	score := 0
	consecutiveCount := 0

	for i, idx := range indices {
		// Base points for match
		score += 10

		// Bonus for consecutive matches
		if i > 0 && indices[i-1] == idx-1 {
			consecutiveCount++
			score += 15
		} else {
			consecutiveCount = 0
		}

		// Bonus for match at start of string
		if idx == 0 {
			score += 10
		}

		// Bonus for match at word boundary
		if idx > 0 {
			prevChar := targetRunes[idx-1]
			if prevChar == '/' || prevChar == '_' || prevChar == '-' || prevChar == ' ' || prevChar == '.' {
				score += 5
			}
		}

		// Bonus for case match
		if targetRunesOriginal[idx] == []rune(pattern)[i] {
			score += 2
		}
	}

	// Penalty for gaps between matches
	if len(indices) > 1 {
		totalGap := indices[len(indices)-1] - indices[0] - (len(indices) - 1)
		score -= totalGap
	}

	// Bonus for shorter target strings (prefer shorter paths)
	score += (1000 - len(targetRunes))

	return score, indices
}

// PerformFuzzyMatch performs fuzzy matching on a list of items and returns sorted matches.
// Items are sorted by score (highest first).
func PerformFuzzyMatch(pattern string, items []string, maxResults int) []FuzzyMatch {
	if pattern == "" {
		// Return all items when no pattern
		var matches []FuzzyMatch
		for _, item := range items {
			if len(matches) >= maxResults {
				break
			}
			matches = append(matches, FuzzyMatch{
				FilePath: item,
				Score:    0,
				Indices:  nil,
			})
		}
		return matches
	}

	var matches []FuzzyMatch

	for _, item := range items {
		score, indices := FuzzyScore(pattern, item)
		if score > 0 {
			matches = append(matches, FuzzyMatch{
				FilePath: item,
				Score:    score,
				Indices:  indices,
			})
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Score > matches[i].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	// Limit results
	if len(matches) > maxResults {
		matches = matches[:maxResults]
	}

	return matches
}

// isWordBoundary checks if a character is a word boundary
func isWordBoundary(r rune) bool {
	return r == '/' || r == '_' || r == '-' || r == ' ' || r == '.' || unicode.IsUpper(r)
}
