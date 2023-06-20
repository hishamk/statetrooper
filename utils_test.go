package statetrooper

import (
	"fmt"
	"testing"
)

type CustomStructStringer struct {
	Name string
	Age  int
}

func (cs CustomStructStringer) String() string {
	return fmt.Sprintf("CustomStruct - Name: %s", cs.Name)
}

type CustomStruct struct {
	Name string
	Age  int
}

func TestStringable(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"Nadia", true}, // String type
		{42, false},     // Non-string type
		{CustomStructStringer{Name: "Yousif"}, true}, // fmt.Stringer type
		{CustomStruct{Name: "Jenna"}, false},         // Non-fmt.Stringer type
	}

	for _, test := range tests {
		actual := stringable(test.input)
		if actual != test.expected {
			t.Errorf("stringable(%v) = %t, expected %t", test.input, actual, test.expected)
		}
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"Nadia", "Nadia"}, // String type
		{42, "42"},         // Non-string type
		{CustomStructStringer{Name: "Yousif"}, "CustomStruct - Name: Yousif"}, // fmt.Stringer type
		{CustomStruct{Name: "Jenna", Age: 12}, "{Jenna 12}"},                  // Non-fmt.Stringer type
	}

	for _, test := range tests {
		actual := toString(test.input)
		if actual != test.expected {
			t.Errorf("toString(%v) = %s, expected %s", test.input, actual, test.expected)
		}
	}
}
