package logging

import (
	"discord-go-music-bot/internal/constants"
	"log"
)

// Simple logging utility with color support for terminal output, used across the bot.

func MessageCreateLog(username, message string) {
	log.Println(constants.ANSIGreen + username + ": " + message + constants.ANSIReset)
}

func InteractionCreateLog(username, command string, args string) {
	log.Println(constants.ANSICyan + username + ": " + command + args + constants.ANSIReset)
}

func ErrorLog(message string) {
	log.Println(constants.ANSIRed + message + constants.ANSIReset)
}

func FatalLog(message string) {
	log.Fatal(constants.ANSIRed + message + constants.ANSIReset)
}

func InfoLog(message string) {
	log.Println(constants.ANSIBlue + message + constants.ANSIReset)
}

func WarningLog(message string) {
	log.Println(constants.ANSIYellow + message + constants.ANSIReset)
}

func DgvoiceLog(message string) {
	log.Println(constants.ANSIRed+ "dgVoice: " + message + constants.ANSIReset)
}
