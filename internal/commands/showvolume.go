package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ShowVolume(s *discordgo.Session, m *discordgo.MessageCreate) {
	state.VolumeMutex.Lock()
	volume, ok := state.Volume[m.GuildID]
	state.VolumeMutex.Unlock()

	if !ok {
		volume = 1.0
	}

	// Convert volume factor back to percentage for display
	volumePercent := volume * 100.0
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current volume is %.1f%%", volumePercent))
}
