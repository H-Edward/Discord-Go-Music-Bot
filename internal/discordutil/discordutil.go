package discordutil

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func GetVoiceConnection(s *discordgo.Session, guildID string) (*discordgo.VoiceConnection, error) {
	vc := s.VoiceConnections[guildID]
	if vc == nil {
		return nil, os.ErrNotExist
	}
	return vc, nil
}

func BotInChannel(s *discordgo.Session, guildID string) bool {
	// determines whether the bot is in the guild's channel
	_, err := GetVoiceConnection(s, guildID)
	return err == nil
}

func JoinUserVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			vc, err := s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, true)
			if err != nil {
				return nil, err
			}
			return vc, nil
		}
	}
	return nil, os.ErrNotExist
}
