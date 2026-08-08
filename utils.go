package gomorphy

import "unicode/utf8"

// Return all splits of a word (taking in account min_reminder and max_prefix_length).
func wordSplits(word string) []prefixWord {
	const minReminder = 3
	const maxPrefixLength = 5
	maxSplit := min(maxPrefixLength, utf8.RuneCountInString(word)-minReminder)
	if maxSplit <= 0 {
		return nil
	}

	result := make([]prefixWord, 0, maxSplit)
	for i := 1; i < 1+maxSplit; i++ {
		leftPart := stringSliceByRuneIndex(word, -1, i)
		rightPart := stringSliceByRuneIndex(word, i, -1)
		result = append(result, prefixWord{prefix: leftPart, unprefixedWord: rightPart})
	}

	return result
}
