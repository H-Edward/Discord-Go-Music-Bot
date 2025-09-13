package handlers

import (
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

func HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Handle slash commands
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	ctx := state.NewInteractionContext(s, i)

	logging.InteractionCreateLog(ctx.User.Username, ctx.CommandName, ctx.ArgumentstoString())
	commandSelector(ctx)
}
