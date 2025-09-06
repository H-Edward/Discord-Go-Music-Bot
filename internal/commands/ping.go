package commands

import "discord-go-music-bot/internal/state"

func Ping(ctx state.Context) {
	ctx.Reply("Ping")
}
