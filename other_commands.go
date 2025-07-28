package main

import (
	"strconv"
	"strings"

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
		"!volume <value> - Sets the volume (0.0 to 2.0)\n" +
		"!currentvolume - Shows the current volume\n" +
		"!nuke <number> - Deletes the specified number of messages\n" +
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
