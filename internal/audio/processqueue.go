package audio

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"
	"fmt"
	"strconv"
	"time"
)

func ProcessQueue(ctx *state.Context) {
	go func() {
		for {
			state.QueueMutex.Lock()
			if len(state.Queue[ctx.GetGuildID()]) == 0 {
				// If the queue is empty, mark the bot as idle and leave the voice channel.
				state.PlayingMutex.Lock()
				state.Playing[ctx.GetGuildID()] = false
				state.PlayingMutex.Unlock()
				state.QueueMutex.Unlock()

				// Wait a moment before disconnecting to avoid rapid connect/disconnect cycles
				time.Sleep(500 * time.Millisecond)

				vc, err := discordutil.GetVoiceConnection(ctx)
				if err == nil {
					vc.Speaking(false)
					vc.Disconnect()
				}
				break
			}

			// Dequeue the next song
			currentURL := state.Queue[ctx.GetGuildID()][0]
			state.Queue[ctx.GetGuildID()] = state.Queue[ctx.GetGuildID()][1:]
			songLength := len(state.Queue[ctx.GetGuildID()])
			state.QueueMutex.Unlock()

			logging.InfoLog("Playing song, " + strconv.Itoa(songLength) + " more in queue ")
			ctx.Reply(fmt.Sprintf("Now playing: %s", currentURL))

			// Create a stop channel for this song
			state.StopMutex.Lock()
			stop := make(chan bool)
			state.StopChannels[ctx.GetGuildID()] = stop
			state.StopMutex.Unlock()

			// Create pause channel
			state.PauseChMutex.Lock()
			pauseCh := make(chan bool, 1) // Buffered channel
			state.PauseChs[ctx.GetGuildID()] = pauseCh
			state.PauseChMutex.Unlock()

			// Initialize pause state
			state.PauseMutex.Lock()
			pauseCh <- state.Paused[ctx.GetGuildID()]
			state.PauseMutex.Unlock()

			done := make(chan bool)
			go playAudio(ctx, currentURL, stop, pauseCh, done)
			<-done

			logging.InfoLog("Song finished, moving to next in queue if available.")

			// Clean up pause channel
			state.PauseChMutex.Lock()
			delete(state.PauseChs, ctx.GetGuildID())
			state.PauseChMutex.Unlock()
		}
	}()
}
