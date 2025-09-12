package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"
)

func CurrentVolume(ctx *state.Context) {
	state.VolumeMutex.Lock()
	defer state.VolumeMutex.Unlock()

	currentVolume, ok := state.Volume[ctx.GetGuildID()]
	if !ok {
		currentVolume = 1.0 // Default volume if not set
		state.Volume[ctx.GetGuildID()] = 1.0
	}
	// Convert to percentage for display
	currentVolume = currentVolume * 100.0 // Convert factor back to percentage
	ctx.Reply(fmt.Sprintf("Current volume is %.1f%%", currentVolume))
}
