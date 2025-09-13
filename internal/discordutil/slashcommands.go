package discordutil

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
	"log"

	"github.com/bwmarrin/discordgo"
)

func SetupSlashCommands(s *discordgo.Session) {
	log.Println(constants.ANSIBlue + "Setting up slash commands" + constants.ANSIReset)
	var theNumberOneAsFloat float64 = 1.0

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
		{Name: "volume", Description: "Set the volume (0-200)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "The volume level (0-200)",
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
					Description: "The number of messages to delete (1-200)",
					Required:    true,
					MaxValue:    200.0,
					MinValue:    &theNumberOneAsFloat,
				},
			},
		},
		{Name: "uptime", Description: "Show the bot's uptime"},
		{Name: "version", Description: "Show the bot's version"},
		{Name: "help", Description: "Show help information"},
		{Name: "oss", Description: "Show the bot's open source information"},
	}

	// Since registering commands takes a bit of time, only register them if they aren't present already
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch existing commands: %v", err)
	}

	for _, cmd := range commands {
		if state.DisabledCommands[cmd.Name] {
			log.Printf(constants.ANSIYellow+"Skipping disabled command: %s"+constants.ANSIReset, cmd.Name)
			continue
		}
		found := false
		for _, existingCmd := range existingCommands {
			if cmd.Name == existingCmd.Name {
				found = true
				break
			}
		}
		if !found {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
			if err != nil {
				log.Fatalf("Could not create '%s' command: %v", cmd.Name, err)
			} else {
				log.Printf(constants.ANSIBlue+"Registered command: %s"+constants.ANSIReset, cmd.Name)
			}
		}
	}
	log.Println(constants.ANSIBlue + "Slash commands setup complete." + constants.ANSIReset)
}
