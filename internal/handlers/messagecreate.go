package handlers

import (
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	logging.MessageCreateLog(m.Author.Username, m.Content)

	if m.Author.Bot || !strings.HasPrefix(m.Content, "!") { // ignore bot messages and messages not starting with '!'
		return
	}

	ctx := state.NewMessageContext(s, m)

	commandSelector(ctx)
}
