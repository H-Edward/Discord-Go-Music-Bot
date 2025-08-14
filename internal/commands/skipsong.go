package commands

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

func SkipSong(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	vc.Speaking(false)

	s.ChannelMessageSend(m.ChannelID, "Skipping current song")

	// The song will stop, and the queue processor will automatically move to the next song
	// We don't need to start a new queue processor
}
