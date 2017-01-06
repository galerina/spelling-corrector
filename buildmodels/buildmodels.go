package main

import (
	"flag"
	"fmt"
	"os"
	"spell-corrector/spell"
)

func main() {
	corpusDir := flag.String("corpusdir", "", "Directory containing corpus files")
	editsDir := flag.String("editsdir", "", "Directory containing edits files")

	flag.Parse()

	if *corpusDir == "" || *editsDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("training corpus: ", *corpusDir)

	lm := spell.NewLanguageModel(*corpusDir)
	fmt.Println("total distinct terms: ", len(lm.Unigrams.WordCounts))
	fmt.Println("total distinct bigrams: ", len(lm.Bigrams.WordCounts))

	lm.Save()
}
