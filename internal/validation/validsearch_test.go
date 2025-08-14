package validation_test

import (
	"discord-go-music-bot/internal/validation"
	"testing"
)

func TestIsValidSearchQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"valid search query", true},
		{"another valid query", true},
		{"", false}, // empty query
		{"a", true}, // single character
		{"this is a very long search query that exceeds the maximum length of two hundred characters, which is not allowed in this test case to ensure that the validation works correctly and does not allow overly long queries to pass through", false}, // too long
		{"invalid@query!", false},                // invalid characters
		{"1234567890", true},                     // numeric query
		{"special characters !@#$%^&*()", false}, // special characters
		{"!", false},
		{"@", false},
		{"#", false},
		{"$", false},
		{"%", false},
		{"^", false},
		{"&", false},
		{"*", false},
		{"(", false},
		{")", false},
		{"\\", false},
		{"\"", false},
		{"'", false},
		{";", false},
		{"<", false},
		{">", false},
		{"?", false},
		{"[", false},
		{"]", false},
		{"{", false},
		{"}", false},
		{"|", false},
		{"`", false},
		{"~", false},
	}

	for _, test := range tests {
		result := validation.IsValidSearchQuery(test.input)
		if result != test.expected {
			t.Errorf("isValidSearchQuery(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
