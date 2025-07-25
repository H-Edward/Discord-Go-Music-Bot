package main

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", true},
		{"http://youtu.be/dQw4w9WgXcQ", true},
		{"invalid-url", false},
		{"ftp://example.com", false},
		{"https://malicious-site.com", false},

		// Valid YouTube URLs - watch
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", true},
		{"http://youtube.com/watch?v=abcdefghijk", true},
		{"https://m.youtube.com/watch?v=ABCDEFGHIJK", true},
		{"youtube.com/watch?v=12345678901", true},
		{"www.youtube.com/watch?v=a1b2c3d4e5F", true},
		{"youtube.com/watch?v=dQw4w9WgXcQ&list=PL12345", true},

		// Valid YouTube URLs - shorts
		{"https://www.youtube.com/shorts/dQw4w9WgXcQ", true},
		{"youtube.com/shorts/ABCDEFGHIJK", true},
		{"m.youtube.com/shorts/12345678901", true},

		// Valid YouTube URLs - embed
		{"https://youtube.com/embed/dQw4w9WgXcQ", true},
		{"www.youtube.com/embed/ABCDEFGHIJK", true},
		{"m.youtube.com/embed/12345678901", true},

		// Valid YouTube URLs - playlist
		{"https://youtube.com/playlist?list=PL1234567890", false},
		{"www.youtube.com/playlist?list=PLabcDEF123", false},
		{"m.youtube.com/playlist?list=PLzyxwvutsr", false},

		// Valid YouTube URLs - youtu.be short links
		{"https://youtu.be/dQw4w9WgXcQ", true},
		{"http://youtu.be/ABCDEFGHIJK", true},
		{"youtu.be/12345678901", true},
		{"www.youtu.be/a1b2c3d4e5F", true},

		// Invalid URLs - wrong domain
		{"https://vimeo.com/123456789", false},
		{"https://youtube.fake.com/watch?v=dQw4w9WgXcQ", false},
		{"https://youtu.be.fake/dQw4w9WgXcQ", false},
		{"https://malicious.youtube.com.evil.com/watch?v=dQw4w9WgXcQ", false},

		// Invalid URLs - wrong path or params
		{"https://youtube.com/watch?foo=bar", false},
		{"https://youtube.com/shorts/12345", false},             // video ID too short
		{"https://youtu.be/12345", false},                       // video ID too short
		{"https://youtube.com/embed/123", false},                // video ID too short
		{"https://youtube.com/playlist?list=", false},           // empty playlist ID
		{"https://youtube.com/playlist?list=INVALID ID", false}, // invalid chars in playlist

		{"youtube.com/watch?v=dQw4w9WgXcQ", true},
		{"youtu.be/dQw4w9WgXcQ", true},

		// Invalid URLs - completely wrong
		{"not a url", false},
		{"", false},
		{"http:///watch?v=dQw4w9WgXcQ", false},
		{"ftp://youtube.com/watch?v=dQw4w9WgXcQ", false},

        // special characters and spaces

        {"you\"tube.com/watch?v=dQw4w9WgXcQ", false},
        {"youtube.com\\/watch?v=dQw4w9WgXcQ", false}, 

	}

	for _, test := range tests {
		result := isValidURL(test.input)
		if result != test.expected {
			t.Errorf("isValidURL(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

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
        {"invalid@query!", false}, // invalid characters
        {"1234567890", true}, // numeric query
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
        result := isValidSearchQuery(test.input)
        if result != test.expected {
            t.Errorf("isValidSearchQuery(%q) = %v; want %v", test.input, result, test.expected)
        }
    }
}