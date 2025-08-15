package audio

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"log"

	"github.com/bwmarrin/discordgo"
)

func playAudio(s *discordgo.Session, guildID, textChannelID, url string, m *discordgo.MessageCreate, stop chan bool, pauseCh chan bool, done chan bool) {
	defer close(done) // Signal when this function exits

	var vc *discordgo.VoiceConnection
	var err error

	if !discordutil.BotInChannel(s, guildID) {
		vc, err = discordutil.JoinUserVoiceChannel(s, m)
		if err != nil {
			log.Println(constants.ANSIRed + "Error joining voice channel: " + err.Error() + constants.ANSIReset)
			s.ChannelMessageSend(textChannelID, "Error joining voice channel.")
			return
		}
	} else {
		vc, err = discordutil.GetVoiceConnection(s, guildID)
		if err != nil {
			log.Println(constants.ANSIRed + "Error getting voice connection: " + err.Error() + constants.ANSIReset)
			s.ChannelMessageSend(textChannelID, "Error with voice connection.")
			return
		}
	}

	songDone := make(chan bool)
	go func() {
		PlayURL(vc, url, stop, pauseCh)
		close(songDone)
	}()

	<-songDone
	log.Println(constants.ANSIBlue + "Song playback complete" + constants.ANSIReset)
}
