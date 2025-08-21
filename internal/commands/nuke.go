package commands

import (
	"discord-go-music-bot/internal/validation"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func NukeMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if the user has permission to manage messages
	if !validation.HasPermission(s, m, discordgo.PermissionManageMessages) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}

	if len(strings.Fields(m.Content)) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !nuke <number of messages>")
		return
	}
	num, err := strconv.Atoi(strings.Fields(m.Content)[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of messages")
		return
	}
	if num < 1 || num > 100 {
		s.ChannelMessageSend(m.ChannelID, "Please specify a number between 1 and 100")
		return
	}
	num++ // Include the command message itself

	messages, err := s.ChannelMessages(m.ChannelID, num, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching messages")
		return
	}
	for _, message := range messages {
		s.ChannelMessageDelete(m.ChannelID, message.ID)
		time.Sleep(20 * time.Millisecond) // Rate limit to avoid hitting Discord's API limits
	}
	s.ChannelMessageSend(m.ChannelID, "Nuked "+strconv.Itoa(num-1)+" messages.")
}
