package commands

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"
	"log"
	"time"
)

func StopSong(ctx state.Context) {
	// Get the voice connection for the guild
	vc, err := discordutil.GetVoiceConnection(ctx)
	if err != nil {
		ctx.Reply("Not in a voice channel")
		return
	}

	// Signal the current song to stop
	state.StopMutex.Lock()
	if stopChan, exists := state.StopChannels[ctx.GetGuildID()]; exists {
		close(stopChan)
		delete(state.StopChannels, ctx.GetGuildID())
	}
	state.StopMutex.Unlock()

	// Clear the queue for the guild
	state.QueueMutex.Lock()
	state.Queue[ctx.GetGuildID()] = []string{}
	state.QueueMutex.Unlock()

	// Mark the bot as not playing
	state.PlayingMutex.Lock()
	state.Playing[ctx.GetGuildID()] = false
	state.PlayingMutex.Unlock()

	// Wait a moment for processes to terminate cleanly
	// then disconnect from the voice channel
	go func() {
		// Give a small delay for processes to clean up
		time.Sleep(500 * time.Millisecond)
		vc.Speaking(false)
		err = vc.Disconnect()
		if err != nil {
			log.Println(constants.ANSIRed + "Error disconnecting from voice channel: " + err.Error() + constants.ANSIReset)
		}
	}()

	// Notify the user
	ctx.Reply("Stopped playback and cleared the queue.")
}
