package commands

import "github.com/bwmarrin/discordgo"

func Oss(s *discordgo.Session, m *discordgo.MessageCreate) {
	repoURL := "https://github.com/H-Edward/Discord-Go-Music-Bot"

	s.ChannelMessageSend(m.ChannelID, "This bot is open source and is available at "+repoURL)

}
