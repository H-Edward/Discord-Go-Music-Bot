package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// These are commands which aren't music-bot related

func nukeMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if the user has permission to manage messages
	if !hasPermission(s, m, discordgo.PermissionManageMessages) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}

	if len(strings.Fields(m.Content)) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !nuke <number of messages>")
		return
	}
	num, err := strconv.Atoi(strings.Fields(m.Content)[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of messages")
		return
	}
	num++ // Include the command message itself

	messages, err := s.ChannelMessages(m.ChannelID, num, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching messages")
		return
	}
	for _, message := range messages {
		s.ChannelMessageDelete(m.ChannelID, message.ID)
	}
	s.ChannelMessageSend(m.ChannelID, "Nuked "+strconv.Itoa(num-1)+" messages.")
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Ping")
}
func pong(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong")
}

func uptime(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !hasPermission(s, m, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}
	timeNow := time.Now()
	uptime := timeNow.Sub(startTime)
	// convert to days, hours, minutes, seconds
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	var uptimeMessage strings.Builder
	if days > 0 {
		if days == 1 {
			uptimeMessage.WriteString("1 day, ")
		} else {
			uptimeMessage.WriteString(strconv.Itoa(days) + " days, ")
		}
	}
	if hours > 0 {
		if hours == 1 {
			uptimeMessage.WriteString("1 hour, ")
		} else {
			uptimeMessage.WriteString(strconv.Itoa(hours) + " hours, ")
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			uptimeMessage.WriteString("1 minute and ")
		} else {

			uptimeMessage.WriteString(strconv.Itoa(minutes) + " minutes and ")
		}
	}
	if seconds > 0 {
		if seconds == 1 {
			uptimeMessage.WriteString("1 second")
		} else {
			uptimeMessage.WriteString(strconv.Itoa(seconds) + " seconds")
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Uptime: "+uptimeMessage.String())
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := "Commands:\n" +
		"!ping - Responds with Pong\n" +
		"!pong - Responds with Ping\n" +
		"!play <url> - Plays a song from the given URL\n" +
		"!search <query> - Searches for a song and plays it\n" +
		"!skip - Skips the current song\n" +
		"!queue - Shows the current queue\n" +
		"!stop - Stops playback and clears the queue\n" +
		"!pause - Pauses playback\n" +
		"!resume - Resumes playback\n" +
		"!volume <value> - Sets the volume (0 to 200)\n" +
		"!currentvolume - Shows the current volume\n" +
		"!nuke <number> - Deletes the specified number of messages\n" +
		"!uptime - Shows how long the bot has been running\n" +
		"!version - Shows a hash-based version of the bot\n" +
		"!help - Shows this help message\n"
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !hasPermission(s, m, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Version: "+GoSourceHash)
}

func unknown(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check .env for how to handle unknown commands
	os.Getenv("UNKNOWN_COMMANDS")
	if os.Getenv("UNKNOWN_COMMANDS") == "ignore" {
		return
	}
	if os.Getenv("UNKNOWN_COMMANDS") == "help" {
		help(s, m)
		return
	}
	if os.Getenv("UNKNOWN_COMMANDS") == "error" {
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Type !help for a list of commands.")
		return
	}
	// if there is no UNKNOWN_COMMANDS in .env, treat as "ignore"
	return
}
