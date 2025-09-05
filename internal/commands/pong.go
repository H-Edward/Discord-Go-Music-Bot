package commands

import "discord-go-music-bot/internal/state"

func Pong(ctx state.Context) {
	ctx.Reply("Pong")
}
