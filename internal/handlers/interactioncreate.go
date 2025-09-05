package handlers

import (
	"discord-go-music-bot/internal/commands"
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
	log.Println(constants.ANSIYellow + ctx.User.Username + ": " + i.Data.Type().String() + constants.ANSIReset)

	switch ctx.CommandName {
	case "ping":
		commands.Pong(*ctx)
	case "pong":
		commands.Ping(*ctx)
	case "play":
		commands.AddSong(*ctx, false) // false as in not a search
	case "search":
		commands.AddSong(*ctx, true) // true as in search for a song
	case "skip":
		commands.SkipSong(*ctx)
	case "queue":
		commands.ShowQueue(*ctx)
	case "stop":
		commands.StopSong(*ctx)
	case "pause", "!resume":
		commands.PauseSong(*ctx)
	case "volume":
		commands.SetVolume(*ctx)
	case "currentvolume":
		commands.CurrentVolume(*ctx)
	case "nuke": // delete n messages
		commands.NukeMessages(*ctx)
	case "uptime":
		commands.Uptime(*ctx)
	case "version":
		commands.Version(*ctx)
	case "help":
		commands.Help(*ctx)
	case "oss":
		commands.Oss(*ctx)
	default:
		commands.Unknown(*ctx)
	}
}
