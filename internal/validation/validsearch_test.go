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

func TestSanitiseSearchQuery(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
		expectedValid  bool
	}{
		{"valid search query", "valid search query", true},
		{"another valid query", "another valid query", true},
		{"", "", false},  // empty query
		{"a", "a", true}, // single character
		{"this is a very long search query that exceeds the maximum length of two hundred characters, which is not allowed in this test case to ensure that the validation works correctly and does not allow overly long queries to pass through", "", false}, // too long
		{"invalid@query!", "invalidquery", true},                                   // invalid characters removed
		{"1234567890", "1234567890", true},                                         // numeric query
		{"special characters !@#$%^&*()", "special characters", true},              // special characters removed
		{"!!!@@@###$$$", "", false},                                                // only special characters
		{"   leading and trailing spaces   ", "leading and trailing spaces", true}, // leading/trailing spaces removed
		{"multiple    spaces", "multiple spaces", true},                            // multiple spaces reduced to single space
		{"valid123 query456", "valid123 query456", true},                           // alphanumeric
		{"mix3d Ch@ract3rs & Spac3s!", "mix3d Chract3rs Spac3s", true},             // mixed valid and invalid characters
		{"!@#$%^&*()_+", "", false},                                                // only special characters
	}

	for _, test := range tests {
		output, valid := validation.SanitiseSearchQuery(test.input)
		if output != test.expectedOutput || valid != test.expectedValid {
			t.Errorf("sanitiseSearchQuery(%q) = (%q, %v); want (%q, %v)", test.input, output, valid, test.expectedOutput, test.expectedValid)
		}
	}
}
