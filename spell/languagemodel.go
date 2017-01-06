package spell

import (
	"bufio"
	"encoding/gob"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type LanguageModel struct {
	Unigrams *Dictionary
	Bigrams  *Dictionary
}

func NewLanguageModel(corpusFilePath string) *LanguageModel {
	lm := LanguageModel{NewDictionary(), NewDictionary()}

	files, err := ioutil.ReadDir(corpusFilePath)
	if err != nil {
		panic(err)
	}

	for _, fi := range files {
		if !fi.IsDir() {
			f, err := os.Open(filepath.Join(corpusFilePath, fi.Name()))
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				tokens := strings.Split(scanner.Text(), " ")
				for i, _ := range tokens {
					lm.Unigrams.Add(tokens[i])

					if i > 0 {
						lm.Bigrams.Add(tokens[i-1] + " " + tokens[i])
					}
				}
			}

			f.Close()
		}
	}

	return &lm
}

func (lm *LanguageModel) unigramProbability(w string) float64 {
	probability := float64(lm.Unigrams.Count(w)) / float64(lm.Unigrams.TermCount)
	if math.IsNaN(probability) || math.IsInf(probability, 0) {
		panic("Invalid probability")
	}

	return probability
}

const lambda = float64(0.1)

// Interpolated bigram probability
// Pint(w2|w1) = λPMLE(w2) + (1 − λ)PMLE(w2|w1)
func (lm *LanguageModel) bigramProbability(w1 string, w2 string) float64 {
	w1_w2 := w1 + " " + w2

	pmle := float64(lm.Bigrams.Count(w1_w2)) / float64(lm.Unigrams.Count(w1))
	probability := lambda*lm.unigramProbability(w2) + (1-lambda)*pmle
	if math.IsNaN(probability) || math.IsInf(probability, 0) {
		panic("Invalid probability")
	}

	return probability
}

// P(w1, w2, ..., wn) = P(w1)P(w2|w1)P(w3|w2)...P(wn|wn−1)
// P(w1) is the unigram probability and all other terms are bigram probabilities
// We calculate log(P(w1, w2, ..., wn)) which is given by
// log(P(w1))+log(P(w2|w1))+log(P(w3|w2)...log(P(wn|wn-1))
func (lm *LanguageModel) LogQueryProbability(query string) float64 {
	terms := strings.Split(query, " ")
	if len(terms) == 0 {
		panic("Invalid query argument")
	}

	probability := math.Log(lm.unigramProbability(terms[0]))
	for i := 1; i < len(terms); i++ {
		probability += math.Log(lm.bigramProbability(terms[i-1], terms[i]))
	}

	return probability
}

func (lm *LanguageModel) TermExists(term string) bool {
	return lm.Unigrams.Count(term) > 0
}

const languageModelSaveFile = "lm"

func (lm *LanguageModel) Save() {
	f, err := os.Create(languageModelSaveFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(lm)
	if err != nil {
		panic(err)
	}
}

func LoadLanguageModel(path string) *LanguageModel {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lm LanguageModel
	dec := gob.NewDecoder(f)
	err = dec.Decode(&lm)
	if err != nil {
		panic(err)
	}

	return &lm
}
