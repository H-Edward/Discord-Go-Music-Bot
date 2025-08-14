package commands

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

func PauseSong(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID

	// Check if the bot is in a voice channel
	if !discordutil.BotInChannel(s, guildID) {
		s.ChannelMessageSend(m.ChannelID, "Not in a voice channel.")
		return
	}

	state.PauseMutex.Lock()
	currentState := state.Paused[guildID]
	state.Paused[guildID] = !currentState // Toggle pause state
	state.PauseMutex.Unlock()

	// Signal the pause channel with the new state
	state.PauseChMutex.Lock()
	if ch, exists := state.PauseChs[guildID]; exists {
		select {
		case ch <- !currentState: // Send the new state
		default: // Channel is full, discard
			// This prevents blocking if the channel is full
		}
	}
	state.PauseChMutex.Unlock()

	if currentState {
		s.ChannelMessageSend(m.ChannelID, "Resumed playback.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Paused playback.")
	}
}
