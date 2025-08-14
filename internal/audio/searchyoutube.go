package audio

import (
	"bytes"
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"
	"log"
	"os/exec"
	"strings"
)

func SearchYoutube(query string) (string, bool) {
	cmd := exec.Command("yt-dlp", "--flat-playlist", "--get-url", "ytsearch1:"+query)
	var outputFromSearch bytes.Buffer
	cmd.Stdout = &outputFromSearch
	err := cmd.Run()
	if err != nil {
		log.Println(state.ANSIRed + "Error: " + err.Error() + state.ANSIReset)
		return "", false
	}

	// Clean up the output - take only the first line
	url := strings.TrimSpace(outputFromSearch.String())
	if idx := strings.Index(url, "\n"); idx > 0 {
		url = url[:idx]
	}

	if url == "" {
		return "", false
	}

	// sanity check with validating the url

	if !validation.IsValidURL(url) {
		return "", false
	}

	return url, true
}
