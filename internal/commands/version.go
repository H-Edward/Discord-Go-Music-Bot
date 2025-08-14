package commands

import (
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"

	"github.com/bwmarrin/discordgo"
)

func Version(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !validation.HasPermission(s, m, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Version: "+state.GoSourceHash)
}
