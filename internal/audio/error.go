package audio

import (
	"discord-go-music-bot/internal/state"
	"log"
)

// OnError gets called by dgvoice when an error is encountered.
var OnError = func(str string, err error) {
	prefix := state.ANSIRed + "dgVoice: " + str

	if err != nil {
		log.Println(prefix + ": " + err.Error() + state.ANSIReset + "\n")
	} else {
		log.Println(prefix + ": Error is nil" + state.ANSIReset + "\n")
	}
}
