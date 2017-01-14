package spell

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func ExampleCorrector() {
	_, filename, _, _ := runtime.Caller(0)
	scriptDir := filepath.Dir(filename)
	lm := LoadLanguageModel(filepath.Join(scriptDir, "../data/lm"))
	corrector1 := NewCorrector(lm, NewUniformEditModel(0.10), 1.0)
	corrector2 := NewCorrector(lm, LoadEmpiricalEditModel(filepath.Join(scriptDir, "../data/em")), 1.0)

	fmt.Println(corrector1.Correct("tast i"))
	fmt.Println(corrector2.Correct("tast i"))
	// Output: past i
	// past i
}
