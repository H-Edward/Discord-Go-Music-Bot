package audio

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ProcessQueue(s *discordgo.Session, guildID, textChannelID string, m *discordgo.MessageCreate) {
	go func() {
		for {
			state.QueueMutex.Lock()
			if len(state.Queue[guildID]) == 0 {
				// If the queue is empty, mark the bot as idle and leave the voice channel.
				state.PlayingMutex.Lock()
				state.Playing[guildID] = false
				state.PlayingMutex.Unlock()
				state.QueueMutex.Unlock()

				// Wait a moment before disconnecting to avoid rapid connect/disconnect cycles
				time.Sleep(500 * time.Millisecond)

				vc, err := discordutil.GetVoiceConnection(s, guildID)
				if err == nil {
					vc.Speaking(false)
					vc.Disconnect()
				}
				break
			}

			// Dequeue the next song
			currentURL := state.Queue[guildID][0]
			state.Queue[guildID] = state.Queue[guildID][1:]
			songLength := len(state.Queue[guildID])
			state.QueueMutex.Unlock()

			log.Printf(state.ANSIBlue+"Playing song, %d more in queue "+state.ANSIReset, songLength)
			s.ChannelMessageSend(textChannelID, fmt.Sprintf("Now playing: %s", currentURL))

			// Create a stop channel for this song
			state.StopMutex.Lock()
			stop := make(chan bool)
			state.StopChannels[guildID] = stop
			state.StopMutex.Unlock()

			// Create pause channel
			state.PauseChMutex.Lock()
			pauseCh := make(chan bool, 1) // Buffered channel
			state.PauseChs[guildID] = pauseCh
			state.PauseChMutex.Unlock()

			// Initialize pause state
			state.PauseMutex.Lock()
			pauseCh <- state.Paused[guildID]
			state.PauseMutex.Unlock()

			done := make(chan bool)
			go playAudio(s, guildID, textChannelID, currentURL, m, stop, pauseCh, done)
			<-done

			log.Println(state.ANSIBlue + "Song finished, moving to next in queue" + state.ANSIReset)

			// Clean up pause channel
			state.PauseChMutex.Lock()
			delete(state.PauseChs, guildID)
			state.PauseChMutex.Unlock()
		}
	}()
}
