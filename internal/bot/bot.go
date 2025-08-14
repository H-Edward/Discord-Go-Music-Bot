package bot

import (
	"discord-go-music-bot/internal/commands"
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
		log.Fatal(state.ANSIRed + "Error loading .env file" + state.ANSIReset)
	}
	state.Token = os.Getenv("DISCORD_BOT_TOKEN")
	if state.Token == "" {
		log.Fatal(state.ANSIRed + "Token not found - check .env file" + state.ANSIReset)
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		log.Fatal(state.ANSIRed + "yt-dlp not found. Please install it with: pip install yt-dlp" + state.ANSIReset)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Fatal(state.ANSIRed + "ffmpeg not found. Please install it with your package manager" + state.ANSIReset)
	}
}

func Run() {
	setup()
	dg, err := discordgo.New("Bot " + state.Token)
	if err != nil {
		log.Fatal(state.ANSIRed + "Error creating Discord session: " + err.Error() + state.ANSIReset)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(state.ANSIRed + "Error opening connection: " + err.Error() + state.ANSIReset)
	}
	defer dg.Close()
	log.Println("Version: " + state.ANSIBold + state.GoSourceHash + state.ANSIReset)
	log.Println(state.ANSIBlue + "Bot is running. Press CTRL-C to exit." + state.ANSIReset)
	select {} // block forever
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Println(state.ANSIYellow + m.Author.Username + ": " + m.Content + state.ANSIReset)

	if m.Author.Bot || !strings.HasPrefix(m.Content, "!") { // ignore bot messages and messages not starting with '!'
		return
	}

	switch strings.Fields(m.Content)[0] {
	case "!ping":
		commands.Pong(s, m)
	case "!pong":
		commands.Ping(s, m)
	case "!play":
		commands.AddSong(s, m, false) // false as in not a search
	case "!search":
		commands.AddSong(s, m, true) // true as in search for a song
	case "!skip":
		commands.SkipSong(s, m)
	case "!queue":
		commands.ShowQueue(s, m)
	case "!stop":
		commands.StopSong(s, m)
	case "!pause", "!resume":
		commands.PauseSong(s, m)
	case "!volume":
		commands.SetVolume(s, m)
	case "!currentvolume":
		commands.CurrentVolume(s, m)
	case "!nuke": // delete n messages
		commands.NukeMessages(s, m)
	case "!uptime":
		commands.Uptime(s, m)
	case "!version":
		commands.Version(s, m)
	case "!help":
		commands.Help(s, m)
	default:
		commands.Unknown(s, m)
	}
}


