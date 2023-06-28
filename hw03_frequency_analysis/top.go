package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const (
	WordSeparators      = "word separators"        // punctuations considered separators
	InWordCharacters    = "in word characters"     // punctuations that are considered part of word
	NotInWordCharacters = "not in word characters" // all characters that are not "InWordCharacters"
)

// Define needed regexp patterns
// WordSeparators are all punctionations without:
//
//	"Connector punctuations" - \p{Pc} - http://www.zuga.net/articles/unicode/category/connector-punctuation/
//	"Dash punctuation" - \p{Pd} - http://www.zuga.net/articles/unicode/category/dash-punctuation/
//
// InWordCharacters - are defined to be "connector punctuations" and "dash punctuations"
// NotInWordCharacters - are defined to contain all characters that are not in InWordCharacters
var expr = map[string]*regexp.Regexp{
	//WordSeparators: regexp.MustCompile(`[.,:;!?()\[\]{}"'\\/#$%&*+=]+`),
	WordSeparators:      regexp.MustCompile(`[^[:^punct:]\p{Pc}\p{Pd}]+`),
	InWordCharacters:    regexp.MustCompile(`[\p{Pd}\p{Pc}]+`),
	NotInWordCharacters: regexp.MustCompile(`[^\p{Pd}\p{Pc}]+`),
}

type WordCount struct {
	Word  string
	Count int
}

// mapToSlice converts a map of strings with int to a slice of WordCount.
func mapToSlice(m map[string]int) []WordCount {
	s := make([]WordCount, 0, len(m))
	for k := range m {
		s = append(s, WordCount{k, m[k]})
	}
	return s
}

// returnFirst10AsSlice returns the first 10 words from the given WordCount slice as a string slice.
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

// returns top 10 word occurences sorted by number of occurences as slice of strings
func Top10(iStr string) []string {
	// create a map to hold the result 3 times smaller then original text
	result := make(map[string]int, len(iStr)/3)

	// Make the incoming text lowercase
	iStr = strings.ToLower(iStr)

	// replace all punctuation marks with spaces
	iStr = expr[WordSeparators].ReplaceAllString(iStr, " ")

	// Split the string into words on whitespace
	str := strings.Fields(iStr)

	// For every separate word
	for _, s := range str {
		// skip any word that is empty and process further those that are not
		if len(s) > 0 {
			// if word has only InWordSeparator chracters (e.g hyphen) and nothing else skip such word
			if expr[InWordCharacters].MatchString(s) && !expr[NotInWordCharacters].MatchString(s) {
				continue
			}
			// increase the count of the word in the map as this word shall be counted
			result[s]++
		}
	}

	// prepare a slice from the map for sorting
	resultSlice := mapToSlice(result)

	// sort the words lexicographically
	sort.Slice(resultSlice, func(i, j int) bool {
		if resultSlice[i].Count == resultSlice[j].Count {
			result := strings.Compare(resultSlice[i].Word, resultSlice[j].Word)
			return result < 0
		}
		return resultSlice[i].Count > resultSlice[j].Count
	})

	// get and return top 10 words
	return returnFirst10AsSlice(resultSlice)
}
