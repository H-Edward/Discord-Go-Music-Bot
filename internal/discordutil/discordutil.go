package discordutil

import (
	"discord-go-music-bot/internal/state"
	"os"

	"github.com/bwmarrin/discordgo"
)

func GetVoiceConnection(ctx *state.Context) (*discordgo.VoiceConnection, error) {
	vc := ctx.GetSession().VoiceConnections[ctx.GetGuildID()]
	if vc == nil {
		return nil, os.ErrNotExist
	}
	return vc, nil
}

func BotInChannel(ctx *state.Context) bool {
	// determines whether the bot is in the guild's channel
	_, err := GetVoiceConnection(ctx)
	return err == nil
}

func JoinUserVoiceChannel(ctx *state.Context) (*discordgo.VoiceConnection, error) {

	guild, err := ctx.GetSession().State.Guild(ctx.GetGuildID())
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.GetUser().ID {
			vc, err := ctx.GetSession().ChannelVoiceJoin(ctx.GuildID, vs.ChannelID, false, true)
			if err != nil {
				return nil, err
			}
			return vc, nil
		}
	}
	return nil, os.ErrNotExist
}

func IsUserInVoiceChannel(ctx *state.Context) bool {
	guild, err := ctx.GetSession().State.Guild(ctx.GetGuildID())
	if err != nil {
		return false
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.GetUser().ID {
			return true
		}
	}
	return false
}
