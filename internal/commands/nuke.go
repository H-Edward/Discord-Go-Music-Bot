package commands

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func NukeMessages(ctx *state.Context) {
	// check if the user has permission to manage messages
	if !validation.HasPermission(ctx, discordgo.PermissionManageMessages) {
		ctx.Reply(constants.EmojiWarning + " You do not have permission to use this command.")
		return
	}

	if ctx.Arguments["count"] == "" {
		ctx.Reply(constants.EmojiInfo + " Usage: !nuke <number of messages>")
		return
	}
	num, err := strconv.Atoi(ctx.Arguments["count"])
	if err != nil {
		ctx.Reply(constants.EmojiWarning + " Invalid number of messages")
		return
	}
	if num < 1 || num > 100 {
		ctx.Reply(constants.EmojiWarning + " Please specify a number between 1 and 100")
		return
	}
	num++ // Include the command message itself

	messages, err := ctx.GetSession().ChannelMessages(ctx.GetChannelID(), num, "", "", "")
	if err != nil {
		ctx.Reply(constants.EmojiWarning + " Error fetching messages")
		return
	}
	for _, message := range messages {
		ctx.GetSession().ChannelMessageDelete(ctx.GetChannelID(), message.ID)
		time.Sleep(20 * time.Millisecond) // Rate limit to avoid hitting Discord's API limits
	}
	ctx.Reply(constants.EmojiSuccess + " Nuked " + strconv.Itoa(num-1) + " messages.")
}
