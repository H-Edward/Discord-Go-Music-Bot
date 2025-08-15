package audio

import (
	"discord-go-music-bot/internal/constants"
	"log"
)

// OnError gets called by dgvoice when an error is encountered.
var OnError = func(str string, err error) {
	prefix := constants.ANSIRed + "dgVoice: " + str

	if err != nil {
		log.Println(prefix + ": " + err.Error() + constants.ANSIReset + "\n")
	} else {
		log.Println(prefix + ": Error is nil" + constants.ANSIReset + "\n")
	}
}
