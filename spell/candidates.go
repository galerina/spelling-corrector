package spell

import (
	"fmt"
	"spell-corrector/stringset"
	"strings"
)

func (lm *LanguageModel) isValid(query string) bool {
	for i := 0; i < len(query)-1; i++ {
		if query[i] == ' ' && query[i+1] == ' ' {
			return false
		}
	}

	for _, s := range strings.Split(query, " ") {
		if !lm.TermExists(s) {
			return false
		}
	}

	return true
}

func transposed(s string, index int) string {
	if index < 0 || index >= len(s)-1 {
		panic(fmt.Sprintf("Invalid argument to transpose(): s: %s, index: %v", s, index))
	}

	return s[:index] + s[index+1:index+2] + s[index:index+1] + s[index+2:]
}

func replaced(s string, char string, index int) string {
	if index < 0 || index >= len(s) {
		panic(fmt.Sprintf("Invalid argument to replaced(): s: %s, char: %s, index: %v", s, char, index))
	}

	return s[:index] + char + s[index+1:]
}

func inserted(s string, char string, index int) string {
	if index < 0 || index > len(s) {
		panic(fmt.Sprintf("Invalid argument to inserted(): s: %s, char: %s, index: %v", s, char, index))
	}

	return s[:index] + char + s[index:]
}

func deleted(s string, index int) string {
	if index < 0 || index >= len(s) {
		panic(fmt.Sprintf("Invalid argument to deleted(): s: %s, index: %v", s, index))
	}

	return s[:index] + s[index+1:]
}

var letters = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
	"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
	"w", "x", "y", "z", " ",
	// "0", "1", "2", "3", "4", "5", "6",
	// "7", "8", "9",
}

func sanitize(query string) string {
	return strings.Join(strings.Split(query, " "), " ")
}

// Edits:
//  Transpose
//  Replacement
//  Insertion
//  Deletion
func (lm *LanguageModel) GetCandidates(query string) *stringset.StringSet {
	query = sanitize(query)

	candidates := stringset.New()
	if lm.isValid(query) {
		candidates.Add(query)
	}

	// Transposes
	for i := 0; i < len(query)-1; i++ {
		s := transposed(query, i)
		if lm.isValid(s) {
			candidates.Add(s)
		}
	}

	// Replacements and deletions
	for i := 0; i < len(query); i++ {
		if query[i] == ' ' {
			continue
		}

		for _, l := range letters {
			s := replaced(query, l, i)
			if lm.isValid(s) {
				candidates.Add(s)
			}
		}

		s := deleted(query, i)
		if lm.isValid(s) {
			candidates.Add(s)
		}
	}

	// Insertion
	for i := 0; i <= len(query); i++ {
		for _, l := range letters {
			s := inserted(query, l, i)
			if lm.isValid(s) {
				candidates.Add(s)
			}
		}
	}

	return candidates
}
