package spell

import (
	"fmt"
)

func ExampleCorrector() {
	lm := LoadLanguageModel("../buildmodels/lm")
	corrector := NewCorrector(lm)

	fmt.Printf(corrector.Correct("tast i"))
	// Output: past i
}
