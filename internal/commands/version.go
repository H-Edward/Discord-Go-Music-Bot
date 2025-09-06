package commands

import (
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"

	"github.com/bwmarrin/discordgo"
)

func Version(ctx state.Context) {
	if !validation.HasPermission(ctx, discordgo.PermissionAdministrator) {
		ctx.Reply("You do not have permission to use this command.")
		return
	}
	ctx.Reply("Version: " + state.GoSourceHash)
}
