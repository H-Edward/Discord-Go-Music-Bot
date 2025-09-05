package handlers

import (
	"discord-go-music-bot/internal/commands"
	"discord-go-music-bot/internal/constants"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Println(constants.ANSIYellow + m.Author.Username + ": " + m.Content + constants.ANSIReset)

	if m.Author.Bot || !strings.HasPrefix(m.Content, "!") { // ignore bot messages and messages not starting with '!'
		return
	}

	switch strings.Fields(m.Content)[0] {
	case "!ping":
		commands.Pong(s, m)
	case "!pong":
		commands.Ping(s, m)
	case "!play":
		commands.AddSong(s, m, false) // false as in not a search
	case "!search":
		commands.AddSong(s, m, true) // true as in search for a song
	case "!skip":
		commands.SkipSong(s, m)
	case "!queue":
		commands.ShowQueue(s, m)
	case "!stop":
		commands.StopSong(s, m)
	case "!pause", "!resume":
		commands.PauseSong(s, m)
	case "!volume":
		commands.SetVolume(s, m)
	case "!currentvolume":
		commands.CurrentVolume(s, m)
	case "!nuke": // delete n messages
		commands.NukeMessages(s, m)
	case "!uptime":
		commands.Uptime(s, m)
	case "!version":
		commands.Version(s, m)
	case "!help":
		commands.Help(s, m)
	case "!oss":
		commands.Oss(s, m)
	default:
		commands.Unknown(s, m)
	}
}
