package audio

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	// Create or retrieve a per-guild Opus encoder to avoid using a single global encoder.
	guildID := v.GuildID

	state.OpusEncoderMutex.Lock()
	enc, exists := state.OpusEncoders[guildID]
	if !exists {
		enc, err = gopus.NewEncoder(constants.FrameRate, constants.Channels, gopus.Audio)
		if err != nil {
			state.OpusEncoderMutex.Unlock()
			OnError("NewEncoder Error", err)
			return
		}
		state.OpusEncoders[guildID] = enc
	}
	state.OpusEncoderMutex.Unlock()

	// Clean up encoder when finished
	defer func() {
		state.OpusEncoderMutex.Lock()
		delete(state.OpusEncoders, guildID)
		state.OpusEncoderMutex.Unlock()
	}()

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return
		}

		// try encoding pcm frame with Opus using the per-guild encoder
		opus, err := enc.Encode(recv, constants.FrameSize, constants.MaxBytes)
		if err != nil {
			OnError("Encoding Error", err)
			return
		}

		if !v.Ready || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
