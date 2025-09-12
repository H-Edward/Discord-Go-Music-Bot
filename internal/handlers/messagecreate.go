package handlers

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Println(constants.ANSIYellow + m.Author.Username + ": " + m.Content + constants.ANSIReset)

	if m.Author.Bot || !strings.HasPrefix(m.Content, "!") { // ignore bot messages and messages not starting with '!'
		return
	}

	ctx := state.NewMessageContext(s, m)

	commandSelector(ctx)
}
