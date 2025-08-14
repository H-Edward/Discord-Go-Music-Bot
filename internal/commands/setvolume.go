package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SetVolume(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !volume <value between 0 and 200>")
		return
	}

	newVolume, err := strconv.ParseFloat(args[1], 64)
	if err != nil || newVolume < 0.0 || newVolume > 200.0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid volume value. Please specify a number between 0 and 200.")
		return
	}
	var preservedVolume float64 = newVolume
	// Normalize the volume to a range of 0.0 to 2.0
	newVolume = newVolume / 100.0 // Convert percentage to a factor

	state.VolumeMutex.Lock()
	if _, ok := state.Volume[m.GuildID]; !ok {
		state.Volume[m.GuildID] = 1.0 // Initialize to default if not set
	}
	state.Volume[m.GuildID] = newVolume
	state.VolumeMutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Volume set to %.2f%%", preservedVolume))
}
