package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"spell-corrector/spell"
)

var config = struct {
	languageModelFilename string
	editModelFilename     string
	dataDirectory         string
}{
	languageModelFilename: "lm",
	editModelFilename:     "em",
	dataDirectory:         "../data",
}

func main() {
	corpusDir := flag.String("corpusdir", "", "Directory containing corpus files")
	editsDir := flag.String("editsdir", "", "Directory containing edits files")

	flag.Parse()

	if *corpusDir == "" || *editsDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Building language model from training corpus: ", *corpusDir)

	lm := spell.NewLanguageModel(*corpusDir)
	fmt.Println("total distinct terms: ", len(lm.Unigrams.WordCounts))
	fmt.Println("total distinct bigrams: ", len(lm.Bigrams.WordCounts))

	fmt.Println("Saving language model...")
	lm.Save(filepath.Join(config.dataDirectory, config.languageModelFilename))

	fmt.Println("Building empirical edit model from edit corpus: ", *editsDir)
	em := spell.NewEmpiricalEditModel(*editsDir)

	fmt.Println("Saving empirical edit model...")
	em.Save(filepath.Join(config.dataDirectory, config.editModelFilename))
}
