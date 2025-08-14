package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func Unknown(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check .env for how to handle unknown commands
	// default case is "ignore"

	unknown_commands := os.Getenv("UNKNOWN_COMMANDS")
	switch unknown_commands {
	case "help":
		Help(s, m)
	case "error":
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Type !help for a list of commands.")
	default:
		return
	}
}
