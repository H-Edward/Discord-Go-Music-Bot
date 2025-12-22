package commands

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
)

func Pong(ctx *state.Context) {
	ctx.Reply(constants.EmojiPing + " Pong")
}
