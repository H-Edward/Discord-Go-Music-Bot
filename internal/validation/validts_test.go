package validation_test

import (
	"discord-go-music-bot/internal/validation"
	"testing"
)

func TestSanitiseTimestamp(t *testing.T) {
	tests := []struct {
		input          string
		expected       bool
		expectedOutput string
	}{
		{"00:00:00", true, "*00:00:00-inf"},
		{"0s", true, "*00:00:00-inf"},
		{"0m", true, "*00:00:00-inf"},
		{"0h", true, "*00:00:00-inf"},
		{"500", true, "*00:08:20-inf"},
		{"1m30s", true, "*00:01:30-inf"},
		{"2h15m", true, "*02:15:00-inf"},
		{"3h45m20s", true, "*03:45:20-inf"},
		{"invalid", false, ""},
		{"1h60m", true, "*02:00:00-inf"},
		{"1h30m60s", true, "*01:31:00-inf"},
		{"-1h", false, ""},
		{"1h-30m", false, ""},
		{"1h30m-20s", false, ""},
		{"1h30m20s extra", false, ""},
		{"1h30m20s!", false, ""},
		{"1h30m20s@", false, ""},
		{"1h30m20s#", false, ""},
		{"1h30m20s$", false, ""},
		{"1h30m20s%", false, ""},
		{"", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sanitised, valid := validation.SanitiseTimestamp(tt.input)
			if valid != tt.expected {
				t.Errorf("SanitiseTimestamp() = %v, want %v", valid, tt.expected)
			}
			if sanitised != tt.expectedOutput {
				t.Errorf("SanitiseTimestamp() = %v, want %v", sanitised, tt.expectedOutput)
			}

			if tt.expected && sanitised == "" {
				t.Errorf("SanitiseTimestamp() returned empty string for valid input %v", tt.input)
			}
		})
	}
}
