package commands

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"
)

func PauseSong(ctx *state.Context) {
	// Regardless of pause or resume, flip the state

	// Check if the bot is in a voice channel
	if !discordutil.BotInChannel(ctx) {
		ctx.Reply("Not in a voice channel.")
		return
	}

	state.PauseMutex.Lock()
	currentState := state.Paused[ctx.GetGuildID()]
	state.Paused[ctx.GetGuildID()] = !currentState // Toggle pause state
	state.PauseMutex.Unlock()

	// Signal the pause channel with the new state
	state.PauseChMutex.Lock()
	if ch, exists := state.PauseChs[ctx.GetGuildID()]; exists {
		select {
		case ch <- !currentState: // Send the new state
		default: // Channel is full, discard
			// This prevents blocking if the channel is full
		}
	}
	state.PauseChMutex.Unlock()

	if currentState {
		ctx.Reply("Resumed playback.")
	} else {
		ctx.Reply("Paused playback.")
	}
}
