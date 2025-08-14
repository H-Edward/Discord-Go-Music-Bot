package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CurrentVolume(s *discordgo.Session, m *discordgo.MessageCreate) {
	state.VolumeMutex.Lock()
	defer state.VolumeMutex.Unlock()

	currentVolume, ok := state.Volume[m.GuildID]
	if !ok {
		currentVolume = 1.0 // Default volume if not set
		state.Volume[m.GuildID] = 1.0
	}
	// Convert to percentage for display
	currentVolume = currentVolume * 100.0 // Convert factor back to percentage
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current volume is %.2f%%", currentVolume))
}
