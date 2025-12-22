package discordutil

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

func SetupSlashCommands(s *discordgo.Session) {
	logging.InfoLog(constants.ANSIBlue + "Setting up slash commands")
	var theNumberOneAsFloat float64 = 1.0

	commands := []*discordgo.ApplicationCommand{
		{Name: "ping", Description: constants.EmojiPing + " Replies with Pong"},
		{Name: "pong", Description: constants.EmojiPing + " Replies with Ping"},
		{Name: "play", Description: constants.EmojiPlay + " Play a Youtube URL",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "The Youtube URL to play",
					Required:    true,
				},
			},
		},
		{Name: "search", Description: constants.EmojiSearch + " Search for a song to play",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The search query",
					Required:    true,
				},
			},
		},
		{Name: "skip", Description: constants.EmojiSkip + " Skip the current song"},
		{Name: "queue", Description: constants.EmojiQueue + " Show the current queue"},
		{Name: "stop", Description: constants.EmojiStop + " Stop playing and clear the queue"},
		{Name: "pause", Description: constants.EmojiPause + " Pause the current song"},
		{Name: "resume", Description: constants.EmojiPlay + " Resume the current song"},
		{Name: "volume", Description: constants.EmojiVolume + " Set the volume (0-200)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "The volume level (0-200)",
					Required:    false, // false since will default to showing current volume
				},
			},
		},
		{Name: "currentvolume", Description: constants.EmojiVolume + " Show the current volume"},
		{Name: "nuke", Description: constants.EmojiNuke + " Delete a number of messages",
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
		{Name: "uptime", Description: constants.EmojiUptime + " Show the bot's uptime"},
		{Name: "version", Description: constants.EmojiInfo + " Show the bot's version"},
		{Name: "help", Description: constants.EmojiInfo + " Show help information"},
		{Name: "oss", Description: constants.EmojiInfo + " Show the bot's open source information"},
	}

	// Since registering commands takes a bit of time, only register them if they aren't present already
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logging.FatalLog("Could not fetch existing commands:" + err.Error())
	}

	for _, cmd := range commands {
		if state.DisabledCommands[cmd.Name] {
			logging.WarningLog("Skipping disabled command:" + cmd.Name)
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
				logging.FatalLog("Could not create command:" + cmd.Name + " " + err.Error())
			} else {
				logging.InfoLog("Registered command: " + cmd.Name)
			}
		}
	}
	logging.InfoLog("Slash commands setup complete.")
}
