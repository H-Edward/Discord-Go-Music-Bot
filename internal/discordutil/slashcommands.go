package discordutil

import (
	"discord-go-music-bot/internal/constants"
	"log"

	"github.com/bwmarrin/discordgo"
)

func SetupSlashCommands(s *discordgo.Session) {

	commands := []*discordgo.ApplicationCommand{
		{Name: "ping", Description: "Replies with Pong"},
		{Name: "pong", Description: "Replies with Ping"},
		{Name: "play", Description: "Play a Youtube URL",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "The Youtube URL to play",
					Required:    true,
				},
			},
		},
		{Name: "search", Description: "Search for a song to play",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The search query",
					Required:    true,
				},
			},
		},
		{Name: "skip", Description: "Skip the current song"},
		{Name: "queue", Description: "Show the current queue"},
		{Name: "stop", Description: "Stop playing and clear the queue"},
		{Name: "pause", Description: "Pause the current song"},
		{Name: "resume", Description: "Resume the current song"},
		{Name: "volume", Description: "Set the volume (0-100)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "The volume level (0-100)",
					Required:    false, // false since will default to showing current volume
				},
			},
		},
		{Name: "currentvolume", Description: "Show the current volume"},
		{Name: "nuke", Description: "Delete a number of messages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "The number of messages to delete (1-100)",
					Required:    true,
				},
			},
		},
		{Name: "uptime", Description: "Show the bot's uptime"},
		{Name: "version", Description: "Show the bot's version"},
		{Name: "help", Description: "Show help information"},
		{Name: "oss", Description: "Show the bot's open source information"},
	}

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Println(constants.ANSIRed + "Cannot create slash command " + v.Name + ": " + err.Error() + constants.ANSIReset)
		}
	}

}
