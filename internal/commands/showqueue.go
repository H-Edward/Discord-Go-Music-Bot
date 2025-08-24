package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ShowQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	state.QueueMutex.Lock()
	defer state.QueueMutex.Unlock()

	if len(state.Queue[m.GuildID]) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Queue is empty.")
		return
	}

	// Make a formatted list of songs, "[N] URL""
	var formattedQueue []string
	for i, song := range state.Queue[m.GuildID] {
		formattedQueue = append(formattedQueue, fmt.Sprintf("[%d] %s", i+1, song))
	}

	s.ChannelMessageSend(m.ChannelID, "Current queue:\n"+strings.Join(formattedQueue, "\n"))
}
