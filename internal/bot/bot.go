package bot

import (
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/handlers"
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func setup() { // find env, get bot token, check dependencies

	if err := godotenv.Load(); err != nil {
		logging.FatalLog("Error loading .env file")
	}
	state.Token = os.Getenv("DISCORD_BOT_TOKEN")
	if state.Token == "" {
		logging.FatalLog("Token not found - check .env file")
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		logging.FatalLog("yt-dlp not found. Please install it with: pip install yt-dlp")
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		logging.FatalLog("ffmpeg not found. Please install it with your package manager")
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
		logging.FatalLog("Error creating Discord session: " + err.Error())
	}

	dg.AddHandler(handlers.HandleMessageCreate)
	dg.AddHandler(handlers.HandleInteractionCreate)

	err = dg.Open()

	discordutil.SetupSlashCommands(dg)

	if err != nil {
		logging.FatalLog("Error opening connection: " + err.Error())
	}
	defer dg.Close()
	logging.InfoLog("Version: " + state.GoSourceHash)
	logging.InfoLog("Bot is running. Press CTRL-C to exit.")
	select {} // block forever
}
