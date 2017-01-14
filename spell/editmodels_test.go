package spell

import "fmt"

func ExampleDeletion() {
	e, _ := getEdit("This", "Thi")
	fmt.Println(e)
	e, _ = getEdit("AThis", "This")
	fmt.Println(e)
	// Output: {is i}
	// {^A ^}
}

func ExampleInsertion() {
	e, _ := getEdit("Thi", "This")
	fmt.Println(e)
	e, _ = getEdit("his", "This")
	fmt.Println(e)
	// Output: {i is}
	// {^ ^T}
}

func ExampleSubstitution() {
	e, _ := getEdit("Thas", "This")
	fmt.Println(e)
	e, _ = getEdit("ahis", "This")
	fmt.Println(e)
	// Output: {a i}
	// {a T}
}

func ExampleTranspose() {
	e, _ := getEdit("hTis", "This")
	fmt.Println(e)
	e, _ = getEdit("Thsi", "This")
	fmt.Println(e)
	// {hT Th}
	// {si is}
}
