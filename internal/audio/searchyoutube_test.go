package audio_test

import (
	"bytes"
	"discord-go-music-bot/internal/audio"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestSearchYoutube(t *testing.T) {
	// The following tests will be run only if yt-dlp is available in the environment.
	// These tests are more for testing whether the software will actually work
	// since yt-dlp is quite buggy and can break easily, so we want to catch that as early as possible.
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp not found in PATH, skipping SearchYoutube tests")
	}

	// Simple test for "Never Gonna Give You Up
	url, found := audio.SearchYoutube("Never Gonna Give You Up")
	if !found {
		t.Error("Expected to find a URL for 'Never Gonna Give You Up', but did not")
	}
	if url == "" {
		t.Error("Expected a non-empty URL for 'Never Gonna Give You Up', but got an empty string")
	}

	// As a sanity test, tanother test is merely checking for the word "never" in the title
	// since yt-dlp can mishandle search results and return something completely unrelated.
	command := exec.Command("yt-dlp", "--get-title", url)
	var outputFromTitle bytes.Buffer
	command.Stdout = &outputFromTitle
	err := command.Run()
	if err != nil {
		t.Error("Error getting title from yt-dlp: " + err.Error())
	}

	title := strings.ToLower(outputFromTitle.String())
	if !strings.Contains(title, "never") {
		t.Errorf("Expected title to contain 'never', but got: %s", title)
	}
	fmt.Printf("SearchYoutube test passed, found URL: %s with title: %s\n", url, title)
}
