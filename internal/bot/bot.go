package bot

import (
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/handlers"
	"discord-go-music-bot/internal/state"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func setup() { // find env, get bot token, check dependencies

	if err := godotenv.Load(); err != nil {
		log.Fatal(constants.ANSIRed + "Error loading .env file" + constants.ANSIReset)
	}
	state.Token = os.Getenv("DISCORD_BOT_TOKEN")
	if state.Token == "" {
		log.Fatal(constants.ANSIRed + "Token not found - check .env file" + constants.ANSIReset)
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		log.Fatal(constants.ANSIRed + "yt-dlp not found. Please install it with: pip install yt-dlp" + constants.ANSIReset)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Fatal(constants.ANSIRed + "ffmpeg not found. Please install it with your package manager" + constants.ANSIReset)
	}

	// Parse disabled commands from .env
	disabled := os.Getenv("DISABLED_COMMANDS")
	for _, cmd := range strings.Split(disabled, ",") {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			state.DisabledCommands[cmd] = true
		}
	}
}

func Run() {
	setup()
	dg, err := discordgo.New("Bot " + state.Token)
	if err != nil {
		log.Fatal(constants.ANSIRed + "Error creating Discord session: " + err.Error() + constants.ANSIReset)
	}

	dg.AddHandler(handlers.HandleMessageCreate)
	dg.AddHandler(handlers.HandleInteractionCreate)

	err = dg.Open()

	discordutil.SetupSlashCommands(dg)

	if err != nil {
		log.Fatal(constants.ANSIRed + "Error opening connection: " + err.Error() + constants.ANSIReset)
	}
	defer dg.Close()
	log.Println(constants.ANSIBlue + "Version: " + constants.ANSIBold + state.GoSourceHash + constants.ANSIReset)
	log.Println(constants.ANSIBlue + "Bot is running. Press CTRL-C to exit." + constants.ANSIReset)
	select {} // block forever
}
