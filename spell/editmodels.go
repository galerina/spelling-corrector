package spell

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type EditModel interface {
	EditProbability(correct, edited string) float64
}

type UniformEditModel struct {
	editProbability float64
}

func NewUniformEditModel(editProbability float64) *UniformEditModel {
	return &UniformEditModel{editProbability: editProbability}
}

func (m *UniformEditModel) EditProbability(correct, edited string) float64 {
	return m.editProbability
}

const beginningOfString = "^"

type Edit struct {
	Prev string
	New  string
}

type EmpiricalEditModel struct {
	ProbabilityTable map[Edit]int
	EditsCount       int
	Chargrams        *Dictionary
}

// firstDiferrence returns the index of the first difference between
// the original and edited strings passed as parameters. If the strings
// are identical return -1.
func firstDifference(orig string, edited string) int {
	for i := 0; i < len(orig); i++ {
		if i == len(edited) {
			return i
		} else if orig[i] != edited[i] {
			return i
		}
	}

	if len(orig) < len(edited) {
		return len(orig)
	}

	// Strings are identical
	return -1
}

// getEdit builds and Edit based on the original and edited strings passed as parameters.
// The parameter strings are assumed to be an edit distance of 1 away from each other
// using deletion, insertion, replacement, or transposition
func getEdit(orig string, edited string) (Edit, error) {
	i := firstDifference(orig, edited)

	argError := errors.New(fmt.Sprintf("'%s' and '%s' are not edit distance 1 apart", orig, edited))
	nilEdit := Edit{"", ""}
	if i == -1 {
		return nilEdit, argError
	}

	if len(orig) == len(edited)+1 {
		// Deletion
		if orig[i+1:] != edited[i:] {
			return nilEdit, argError
		}

		if i == 0 {
			return Edit{beginningOfString + orig[i:i+1], beginningOfString}, nil
		} else {
			return Edit{orig[i-1 : i+1], edited[i-1 : i]}, nil
		}
	} else if len(edited) == len(orig)+1 {
		// Insertion
		if edited[i+1:] != orig[i:] {
			return nilEdit, argError
		}

		if i == 0 {
			return Edit{beginningOfString, beginningOfString + edited[i:i+1]}, nil
		} else {
			return Edit{orig[i-1 : i], edited[i-1 : i+1]}, nil
		}
	} else if len(edited) == len(orig) {
		if i+1 < len(edited) && orig[i+1] == edited[i] && orig[i] == edited[i+1] && orig[i+2:] == edited[i+2:] {
			// Transpose
			return Edit{orig[i : i+2], edited[i : i+2]}, nil
		} else if orig[i+1:] == edited[i+1:] {
			// Replacement
			return Edit{orig[i : i+1], edited[i : i+1]}, nil
		}
	}

	return nilEdit, argError
}

func (em *EmpiricalEditModel) Add(e Edit) {
	em.ProbabilityTable[e] = em.ProbabilityTable[e] + 1
	em.EditsCount++
}

func NewEmpiricalEditModel(dir string) *EmpiricalEditModel {
	em := EmpiricalEditModel{map[Edit]int{}, 0, NewDictionary()}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, fi := range files {
		if !fi.IsDir() {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(f)
			uneditedCount := 0
			for scanner.Scan() {
				queryStrings := strings.Split(scanner.Text(), "\t")
				if len(queryStrings) != 2 {
					panic("Invalid file format")
				}
				orig, corrected := queryStrings[0], queryStrings[1]

				for j := 0; j < len(corrected); j++ {
					em.Chargrams.Add(corrected[j : j+1])
					if j == 0 {
						em.Chargrams.Add(beginningOfString)
						em.Chargrams.Add(beginningOfString + corrected[j:j+1])
					}

					if j < len(corrected)-1 {
						em.Chargrams.Add(corrected[j : j+2])
					}
				}

				if orig != corrected {
					edit, err := getEdit(corrected, orig)
					if err != nil {
						fmt.Println(err)
					} else {
						em.Add(edit)
					}
				} else {
					uneditedCount++
				}
			}

			// fmt.Printf("Unedited count: %v\n", uneditedCount)
			// fmt.Printf("Edits count: %v\n", em.EditsCount)

			f.Close()
		}
	}

	return &em
}

func (em *EmpiricalEditModel) Save(path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(em)
	if err != nil {
		panic(err)
	}
}

func LoadEmpiricalEditModel(path string) *EmpiricalEditModel {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var em EmpiricalEditModel
	dec := gob.NewDecoder(f)
	err = dec.Decode(&em)
	if err != nil {
		panic(err)
	}

	return &em
}

const alphabetSize = 26

func (m *EmpiricalEditModel) EditProbability(orig, edit string) float64 {
	e, _ := getEdit(orig, edit)

	if m.Chargrams.Count(e.Prev) == 0 {
		return 0
	} else {
		smoothingDenominator := 0.5 * float64(alphabetSize)
		return (float64(m.ProbabilityTable[e]) + 0.5) / (float64(m.Chargrams.Count(e.Prev)) + smoothingDenominator)
	}
}
