package stringset

import (
	"fmt"
	"testing"
)

func TestStringset(t *testing.T) {
	s := New()
	s.Add("this")
	s.Add("this")
	s.Add("that")

	fmt.Println(s.Contains("this"))
	fmt.Println(s.Contains("that"))
	fmt.Println(s.Contains("there"))
	// Output: true
	// Output: true
	// Output: false
}
