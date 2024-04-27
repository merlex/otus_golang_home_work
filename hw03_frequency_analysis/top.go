package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(t string) []string {
	wordsCount := make(map[string]int)

	for _, word := range strings.Fields(strings.ToLower(t)) {
		wordsCount[strings.Trim(word, ",.!?:;")]++
	}

	words := make([]string, 0, len(wordsCount))
	for key := range wordsCount {
		if key != "" && key != "-" {
			words = append(words, key)
		}
	}

	sort.Slice(words, func(i, j int) bool {
		if wordsCount[words[i]] == wordsCount[words[j]] {
			return words[i] < words[j]
		}
		return wordsCount[words[i]] > wordsCount[words[j]]
	})

	return words[:min(len(words), 10)]
}
