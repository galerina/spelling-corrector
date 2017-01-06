package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"spell-corrector/spell"
)

const (
	dataDir       = "../data/dev_set"
	queryFilename = "queries.txt"
	goldFilename  = "gold.txt"
)

// Baseline (no correction): .479 accuracy
// Google (google.txt): .831 accuracy
func main() {
	queryFile, err := os.Open(filepath.Join(dataDir, queryFilename))
	if err != nil {
		panic(err)
	}
	defer queryFile.Close()

	goldFile, err := os.Open(filepath.Join(dataDir, goldFilename))
	if err != nil {
		panic(err)
	}
	defer goldFile.Close()

	queryScanner := bufio.NewScanner(queryFile)
	goldScanner := bufio.NewScanner(goldFile)

	lm := spell.LoadLanguageModel("../buildmodels/lm")
	corrector := spell.NewCorrector(lm)

	correct := 0
	total := 0
	for queryScanner.Scan() && goldScanner.Scan() {
		correctedQuery := corrector.Correct(queryScanner.Text())
		fmt.Println("Corrected:")
		fmt.Println(queryScanner.Text())
		fmt.Println("To:")
		fmt.Println(correctedQuery)
		if correctedQuery == goldScanner.Text() {
			correct++
		}
		total++
	}

	fmt.Printf("Correct: %v, Total: %v, Accuracy: %v\n", correct, total, float32(correct)/float32(total))
}
