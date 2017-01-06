package spell

import "testing"

type terms []struct {
	term  string
	count int
}

type termTest struct {
	t         terms
	termCount int
}

var test1 = termTest{
	t: terms{
		{"love", 1},
		{"hate", 1},
		{"feeling", 7},
	},
	termCount: 9,
}

func checkExpectedAgainstActual(expected int, actual int, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected=%d does not equal actual=%d", expected, actual)
	}
}

func TestDictionary(t *testing.T) {
	d := NewDictionary()

	for _, tc := range test1.t {

		for i := 0; i < tc.count; i++ {
			d.Add(tc.term)
		}
	}

	checkExpectedAgainstActual(test1.termCount, d.TermCount, t)

	for _, tc := range test1.t {
		checkExpectedAgainstActual(tc.count, d.Count(tc.term), t)
	}
}
