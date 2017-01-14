package spell

import (
	"fmt"
	"testing"
)

func ExampleTransposed() {
	fmt.Println(transposed("act", 0))
	fmt.Println(transposed("abcd", 0))
	fmt.Println(transposed("jokester", 6))
	// Output:
	// cat
	// bacd
	// jokestre
}

func ExampleReplaced() {
	fmt.Println(replaced("abcd", "k", 0))
	fmt.Println(replaced("killjoy", "z", 6))
	// Output:
	// kbcd
	// killjoz
}

func ExampleInserted() {
	fmt.Println(inserted("abcd", "j", 0))
	fmt.Println(inserted("cameron", "z", 7))
	// Output:
	// jabcd
	// cameronz
}

func ExampleDeleted() {
	fmt.Println(deleted("abcd", 0))
	fmt.Println(deleted("christmas", 8))
	// Output:
	// bcd
	// christma
}

var tests = []struct {
	query      string
	candidates []string
}{
	{query: "act", candidates: []string{"cat", "ace", "tact"}},
	{query: "snooop", candidates: []string{"snoop"}},
}

func TestCandidates(t *testing.T) {
	lm := LoadLanguageModel("../buildmodels/lm")

	for _, test := range tests {
		c := lm.GetCandidates(test.query)
		for _, expected := range test.candidates {
			if !c.Contains(expected) {
				t.Errorf("Query: %s, expected value %s not in candidates\n", test.query, expected)
			}
		}

		/*
			fmt.Printf("Query = %s\n", test.query)
			fmt.Printf("%10s%10s\n", "Edit", "ln prob")
			for cand := range *c {
				fmt.Printf("%10s%10.5v\n", cand, lm.LogQueryProbability(cand))
			}
		*/
	}
}
