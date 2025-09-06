package commands

import "discord-go-music-bot/internal/state"

func Oss(ctx state.Context) {
	repoURL := "https://github.com/H-Edward/Discord-Go-Music-Bot"

	ctx.Reply("This bot is open source and is available at " + repoURL)

}
