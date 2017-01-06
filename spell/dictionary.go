package spell

type Dictionary struct {
	TermCount  int
	WordCounts map[string]int
}

func NewDictionary() *Dictionary {
	d := Dictionary{0, map[string]int{}}
	return &d
}

func (d *Dictionary) Add(term string) {
	d.TermCount++

	d.WordCounts[term] = d.WordCounts[term] + 1
}

func (d *Dictionary) Count(term string) int {
	return d.WordCounts[term]
}
