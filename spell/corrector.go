package spell

import (
	"math"
	"spell-corrector/stringset"
)

type Corrector struct {
	lm *LanguageModel
	em EditModel
	mu float64
}

func NewCorrector(lm *LanguageModel, em EditModel, mu float64) *Corrector {
	return &Corrector{lm: lm, em: em, mu: mu}
}

// max returns the max string using the less function passed in as a parameter
func max(set *stringset.StringSet, less func(s1 string, s2 string) bool) string {
	maxString := set.GetAny()
	for s := range *set {
		if less(maxString, s) {
			maxString = s
		}
	}

	return maxString
}

func (corrector *Corrector) queryProb(s string) float64 {
	return math.Exp(corrector.lm.LogQueryProbability(s))
}

func (corrector *Corrector) editProbability(orig string, edited string, distance int) float64 {
	if distance == 0 {
		return 0.90
	} else {
		return math.Pow(corrector.em.EditProbability(orig, edited), float64(distance))
	}
}

func distance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	} else {
		return 1
	}
}

func (corrector *Corrector) Correct(query string) string {
	candidates := corrector.lm.GetCandidates(query)
	return max(candidates, func(s1 string, s2 string) bool {
		dist := distance(s1, query)
		s1Probability := corrector.editProbability(s1, query, dist) *
			math.Pow(corrector.queryProb(s1), corrector.mu)
		dist = distance(s2, query)
		s2Probability := corrector.editProbability(s2, query, dist) *
			math.Pow(corrector.queryProb(s2), corrector.mu)
		return s1Probability < s2Probability
	})
}
