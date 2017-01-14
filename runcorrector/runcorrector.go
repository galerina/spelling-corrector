package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"spell-corrector/spell"
)

const (
	dataDir           = "../data/dev_set"
	queryFilename     = "queries.txt"
	goldFilename      = "gold.txt"
	languageModelFile = "../data/lm"
	editsModelFile    = "../data/em"
)

// Baseline (no correction): .479 accuracy
// Google (google.txt): .831 accuracy
// Uniform edit model with distance 1 edits and language model: 0.7956044
// Initial empirical edit model w/ language model: 0.8175824
func main() {
	editModel := flag.String("editmodel", "uniform", "<uniform|empirical>, uniform is default")
	mu := flag.Float64("mu", 1.0, "Weight for the edit model probability")

	flag.Parse()

	if *editModel != "uniform" && *editModel != "empirical" {
		flag.PrintDefaults()
		os.Exit(1)
	}

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

	lm := spell.LoadLanguageModel(languageModelFile)
	var em spell.EditModel
	if *editModel == "uniform" {
		em = spell.NewUniformEditModel(0.10)
	} else if *editModel == "empirical" {
		em = spell.LoadEmpiricalEditModel(editsModelFile)
	} else {
		panic(fmt.Sprintf("Unexpected edit model type: %s", *editModel))
	}

	corrector := spell.NewCorrector(lm, em, *mu)

	correct := 0
	total := 0
	for queryScanner.Scan() && goldScanner.Scan() {
		correctedQuery := corrector.Correct(queryScanner.Text())
		if correctedQuery == goldScanner.Text() {
			correct++
		} else {
			fmt.Printf("Incorrectly corrected: '%s'\n", queryScanner.Text())
			fmt.Printf("To: '%s'\n", correctedQuery)
			fmt.Printf("Should be: '%s'\n", goldScanner.Text())
		}
		total++
	}

	fmt.Printf("Correct: %v, Total: %v, Accuracy: %v\n", correct, total, float32(correct)/float32(total))
}
