package audio

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

func playAudio(ctx *state.Context, url string, stop chan bool, pauseCh chan bool, done chan bool) {
	defer close(done) // Signal when this function exits

	var vc *discordgo.VoiceConnection
	var err error

	if !discordutil.BotInChannel(ctx) {
		vc, err = discordutil.JoinUserVoiceChannel(ctx)
		if err != nil {
			logging.ErrorLog("Error joining voice channel: " + err.Error())
			ctx.Reply("Error joining voice channel.")
			return
		}
	} else {
		vc, err = discordutil.GetVoiceConnection(ctx)
		if err != nil {
			logging.ErrorLog("Error getting voice connection: " + err.Error())
			ctx.Reply("Error with voice connection.")
			return
		}
	}

	songDone := make(chan bool)
	go func() {
		PlayURL(vc, url, stop, pauseCh)
		close(songDone)
	}()

	<-songDone
	logging.InfoLog("Song playback complete")
}
