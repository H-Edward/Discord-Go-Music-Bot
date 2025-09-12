package handlers

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
	"log"

	"github.com/bwmarrin/discordgo"
)

func HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Handle slash commands
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	ctx := state.NewInteractionContext(s, i)
	//	log.Println(constants.ANSIYellow + m.Author.Username + ": " + m.Content + constants.ANSIReset)

	log.Println(constants.ANSICyan + ctx.User.Username + ": " + ctx.CommandName + ctx.ArgumentstoString() + constants.ANSIReset)

	commandSelector(ctx)
}
