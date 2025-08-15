package commands

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func StopSong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get the voice connection for the guild
	vc, err := discordutil.GetVoiceConnection(s, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not in a voice channel")
		return
	}

	// Signal the current song to stop
	state.StopMutex.Lock()
	if stopChan, exists := state.StopChannels[m.GuildID]; exists {
		close(stopChan)
		delete(state.StopChannels, m.GuildID)
	}
	state.StopMutex.Unlock()

	// Clear the queue for the guild
	state.QueueMutex.Lock()
	state.Queue[m.GuildID] = []string{}
	state.QueueMutex.Unlock()

	// Mark the bot as not playing
	state.PlayingMutex.Lock()
	state.Playing[m.GuildID] = false
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
	s.ChannelMessageSend(m.ChannelID, "Stopped playback and cleared the queue.")
}
