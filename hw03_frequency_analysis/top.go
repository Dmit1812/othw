package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// Define needed patterns
// var expr = map[string]*regexp.Regexp{
// 	"punctuation": regexp.MustCompile("[[:punct:]]+"),
// 	"hyphen": regexp.MustCompile("-"),
// }

type WordCount struct {
	Word  string
	Count int
}

func mapToSlice(m map[string]int) []WordCount {
	s := make([]WordCount, 0, len(m))
	for k := range m {
		s = append(s, WordCount{k, m[k]})
		//fmt.Printf("Word: %s, Count: %d\n", k, m[k])
	}
	return s
}

func returnFirst10AsSlice(s []WordCount) []string {
	result := make([]string, 0, len(s))
	for _, w := range s {
		result = append(result, w.Word)
	}
	i := len(result)
	if i > 10 {
		i = 10
	}
	return result[:i]
}

func Top10(iStr string) []string {
	result := make(map[string]int, len(iStr)/3)

	str := strings.Fields(iStr)
	for _, s := range str {
		if len(s) > 0 {
			result[s]++
		}
	}

	resultSlice := mapToSlice(result)

	sort.Slice(resultSlice, func(i, j int) bool {
		if resultSlice[i].Count == resultSlice[j].Count {
			return resultSlice[i].Word < resultSlice[j].Word
		}
		return resultSlice[i].Count > resultSlice[j].Count
	})

	// Extract word according to rules and add a count to the map
	// Set 1
	// Rule 1 - case sensitive, words are separated by spaces
	// Rule 2 - punctuation marks are parts of words
	// Rule 3 - "-" is a separate word

	// Set 2 - override rules of Set 1
	// Rule 4 - case insensitive, words are separated by spaces and punctuation marks
	// Rule 5 - punctuation marks are separators like space, not counted
	// Rule 6 - "-" is a letter from the word
	// Place your code here.

	// get top 10 words

	// sort the words lexicographically

	return returnFirst10AsSlice(resultSlice)
}
