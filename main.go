package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data

	ffmpegBufferSize int = 16384 // 2^14

	maxClampValue float64 = 32767  // max volume
	minClampValue float64 = -32768 // min volume

	// ANSI color codes for terminal output
	ANSIBold   = "\033[1m"
	ANSIBlue   = "\033[34m"
	ANSIYellow = "\033[33m"
	ANSIRed    = "\033[31m"
	ANSIReset  = "\033[0m"
)

var (
	token        string
	queue        = make(map[string][]string) // Guild ID -> Queue of URLs
	queueMutex   sync.Mutex
	playing      = make(map[string]bool)
	playingMutex sync.Mutex
	paused       = make(map[string]bool) // Guild ID -> Pause state
	pauseMutex   sync.Mutex
	volume       = make(map[string]float64) // Guild ID -> Volume
	volumeMutex  sync.Mutex
	opusEncoder  *gopus.Encoder
	stopChannels = make(map[string]chan bool)
	stopMutex    sync.Mutex
	pauseChs     = make(map[string]chan bool) // Map of guild ID to pause channels
	pauseChMutex sync.Mutex

	GoSourceHash string // short hash of all go source files
)

func setup() { // find env, get bot token

	if err := godotenv.Load(); err != nil {
		log.Fatal(ANSIRed + "Error loading .env file" + ANSIReset)
	}
	token = os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal(ANSIRed + "Token not found - check .env file" + ANSIReset)
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		log.Fatal(ANSIRed + "yt-dlp not found. Please install it with: pip install yt-dlp" + ANSIReset)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Fatal(ANSIRed + "ffmpeg not found. Please install it with your package manager" + ANSIReset)
	}
}

func main() {
	setup()
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(ANSIRed + "Error creating Discord session: " + err.Error() + ANSIReset)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(ANSIRed + "Error opening connection: " + err.Error() + ANSIReset)
	}
	defer dg.Close()
	log.Println("Version: " + ANSIBold + GoSourceHash + ANSIReset)
	log.Println(ANSIBlue + "Bot is running. Press CTRL-C to exit." + ANSIReset)
	select {} // block forever
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Println(ANSIYellow + m.Author.Username + ": " + m.Content + ANSIReset)

	if m.Author.Bot { // ignore bot messages
		return
	}
	switch strings.Fields(m.Content)[0] {
	case "!ping":
		pong(s, m)
	case "!pong":
		ping(s, m)
	case "!play":
		// make sure its length is suitable
		if len(m.Content) < 6 {
			s.ChannelMessageSend(m.ChannelID, "Invalid URL")
			return
		}

		if !(isValidURL(m.Content[6:]) || m.Content[6:] == "") {

			s.ChannelMessageSend(m.ChannelID, "Invalid URL")
			return
		}
		addSong(s, m, false) // false as in not a search
	case "!search":
		// search for a song on youtube
		// add the first result to the queue
		if len(m.Content) < 7 {
			s.ChannelMessageSend(m.ChannelID, "Invalid search")
			return
		}

		addSong(s, m, true)

	case "!skip":
		skipSong(s, m)
	case "!queue":
		showQueue(s, m)
	case "!stop":
		stopSong(s, m)
	case "!pause", "!resume":
		pauseSong(s, m)
	case "!volume":
		setVolume(s, m)
	case "!currentvolume":
		currentVolume(s, m)
	case "!nuke": // delete n messages
		nukeMessages(s, m)
	case "!version":
		version(s, m)
	case "!help":
		help(s, m)
	default:
		return // ignore other messages
	}
}

func addSong(s *discordgo.Session, m *discordgo.MessageCreate, mode bool) { // mode (false for play, true for search)

	if mode {
		searchQuery := strings.TrimSpace(m.Content[8:])

		if !isValidSearchQuery(searchQuery) {
			s.ChannelMessageSend(m.ChannelID, "Invalid search query")
			return
		}

		cmd := exec.Command("yt-dlp", "--flat-playlist", "--get-url", "ytsearch1:"+searchQuery)
		var outputFromSearch bytes.Buffer
		cmd.Stdout = &outputFromSearch
		err := cmd.Run()
		if err != nil {
			log.Println(ANSIRed + "Error: " + err.Error() + ANSIReset)
			s.ChannelMessageSend(m.ChannelID, "Error while searching")
			return
		}

		// Clean up the output - take only the first line
		url := strings.TrimSpace(outputFromSearch.String())
		if idx := strings.Index(url, "\n"); idx > 0 {
			url = url[:idx]
		}

		if url == "" {
			s.ChannelMessageSend(m.ChannelID, "No results found")
			return
		}

		// sanity check with validating the url

		if !isValidURL(url) {
			s.ChannelMessageSend(m.ChannelID, "Error with found URL")
		}

		s.ChannelMessageSend(m.ChannelID, "Found: "+url)

		queueMutex.Lock()

		queue[m.GuildID] = append(queue[m.GuildID], url)
	} else {
		playURL := strings.TrimSpace(m.Content[6:])

		if !isValidURL(playURL) {
			s.ChannelMessageSend(m.ChannelID, "Invalid or disallowed URL")
			return
		}
		queueMutex.Lock()

		queue[m.GuildID] = append(queue[m.GuildID], m.Content[6:])
	}
	isAlreadyPlaying := playing[m.GuildID]
	queueMutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, "Added to queue.")

	if !isAlreadyPlaying {
		// Start processing the queue if the bot is idle
		playingMutex.Lock()
		playing[m.GuildID] = true
		playingMutex.Unlock()
		processQueue(s, m.GuildID, m.ChannelID, m)
	}
}

func processQueue(s *discordgo.Session, guildID, textChannelID string, m *discordgo.MessageCreate) {
	go func() {
		for {
			queueMutex.Lock()
			if len(queue[guildID]) == 0 {
				// If the queue is empty, mark the bot as idle and leave the voice channel.
				playingMutex.Lock()
				playing[guildID] = false
				playingMutex.Unlock()
				queueMutex.Unlock()

				// Wait a moment before disconnecting to avoid rapid connect/disconnect cycles
				time.Sleep(500 * time.Millisecond)

				vc, err := getVoiceConnection(s, guildID)
				if err == nil {
					vc.Speaking(false)
					vc.Disconnect()
				}
				break
			}

			// Dequeue the next song
			currentURL := queue[guildID][0]
			queue[guildID] = queue[guildID][1:]
			songLength := len(queue[guildID])
			queueMutex.Unlock()

			log.Printf(ANSIBlue+"Playing song, %d more in queue "+ANSIReset, songLength)
			s.ChannelMessageSend(textChannelID, fmt.Sprintf("Now playing: %s", currentURL))

			// Create a stop channel for this song
			stopMutex.Lock()
			stop := make(chan bool)
			stopChannels[guildID] = stop
			stopMutex.Unlock()

			// Create pause channel
			pauseChMutex.Lock()
			pauseCh := make(chan bool, 1) // Buffered channel
			pauseChs[guildID] = pauseCh
			pauseChMutex.Unlock()

			// Initialize pause state
			pauseMutex.Lock()
			pauseCh <- paused[guildID]
			pauseMutex.Unlock()

			done := make(chan bool)
			go playAudio(s, guildID, textChannelID, currentURL, m, stop, pauseCh, done)
			<-done

			log.Println(ANSIBlue + "Song finished, moving to next in queue" + ANSIReset)

			// Clean up pause channel
			pauseChMutex.Lock()
			delete(pauseChs, guildID)
			pauseChMutex.Unlock()
		}
	}()
}

func playAudio(s *discordgo.Session, guildID, textChannelID, url string, m *discordgo.MessageCreate, stop chan bool, pauseCh chan bool, done chan bool) {
	defer close(done) // Signal when this function exits

	var vc *discordgo.VoiceConnection
	var err error

	if !botInChannel(s, guildID) {
		vc, err = joinUserVoiceChannel(s, m)
		if err != nil {
			log.Println(ANSIRed + "Error joining voice channel: " + err.Error() + ANSIReset)
			s.ChannelMessageSend(textChannelID, "Error joining voice channel.")
			return
		}
	} else {
		vc, err = getVoiceConnection(s, guildID)
		if err != nil {
			log.Println(ANSIRed + "Error getting voice connection: " + err.Error() + ANSIReset)
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
	log.Println(ANSIBlue + "Song playback complete" + ANSIReset)
}

func getVoiceConnection(s *discordgo.Session, guildID string) (*discordgo.VoiceConnection, error) {
	vc := s.VoiceConnections[guildID]
	if vc == nil {
		return nil, os.ErrNotExist
	}
	return vc, nil
}

func botInChannel(s *discordgo.Session, guildID string) bool {
	// determines whether the bot is in the guild's channel
	_, err := getVoiceConnection(s, guildID)
	return err == nil
}

func joinUserVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			vc, err := s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, true)
			if err != nil {
				return nil, err
			}
			return vc, nil
		}
	}
	return nil, os.ErrNotExist
}

// OnError gets called by dgvoice when an error is encountered.
var OnError = func(str string, err error) {
	prefix := ANSIRed + "dgVoice: " + str

	if err != nil {
		log.Println(prefix + ": " + err.Error() + ANSIReset + "\n")
	} else {
		log.Println(prefix + ": Error is nil" + ANSIReset + "\n")
	}
}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		OnError("NewEncoder Error", err)
		return
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			OnError("PCM Channel closed", nil)
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
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

// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayURL(v *discordgo.VoiceConnection, url string, stop <-chan bool, pauseCh <-chan bool) {

	if !isValidURL(url) {
		OnError("Invalid URL", nil)
		return
	}

	// Create the yt-dlp command to download only the audio with best quality
	ytDlpCmd := exec.Command("yt-dlp",
		"-f", "bestaudio",
		"--no-playlist",
		"-o", "-",
		url) // Get only audio, best quality

	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	// Setup proper cleanup to ensure processes terminate
	defer func() {
		if ytDlpCmd.Process != nil {
			ytDlpCmd.Process.Kill()
			ytDlpCmd.Wait()
		}
		if ffmpegCmd.Process != nil {
			ffmpegCmd.Process.Kill()
			ffmpegCmd.Wait()
		}
	}()

	// Connect yt-dlp output to ffmpeg input
	ytDlpOut, err := ytDlpCmd.StdoutPipe()
	if err != nil {
		OnError("yt-dlp StdoutPipe Error", err)
		return
	}

	ffmpegIn, err := ffmpegCmd.StdinPipe()
	if err != nil {
		OnError("ffmpeg StdinPipe Error", err)
		return
	}

	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		OnError("ffmpeg StdoutPipe Error", err)
		return
	}

	// Start the yt-dlp and ffmpeg processes
	err = ytDlpCmd.Start()
	if err != nil {
		OnError("yt-dlp Start Error", err)
		return
	}

	err = ffmpegCmd.Start()
	if err != nil {
		OnError("ffmpeg Start Error", err)
		return
	}

	// Pipe yt-dlp output to ffmpeg input
	go func() {
		_, err := io.Copy(ffmpegIn, ytDlpOut)
		if err != nil {
			OnError("Error copying yt-dlp output to ffmpeg input", err)
		}
		ffmpegIn.Close() // Important: close the pipe when done
	}()

	// Set up reading from ffmpeg output
	ffmpegbuf := bufio.NewReaderSize(ffmpegOut, ffmpegBufferSize)

	// Handle stopping ffmpeg process if needed
	go func() {
		select {
		case <-stop:
			err := ytDlpCmd.Process.Kill()
			if err != nil {
				OnError("Error killing yt-dlp process", err)
			}
			err = ffmpegCmd.Process.Kill()
			if err != nil {
				OnError("Error killing ffmpeg process", err)
			}
		}
	}()

	// Make sure voice connection is ready before starting
	time.Sleep(100 * time.Millisecond)

	// Set voice speaking status
	err = v.Speaking(true)
	if err != nil {
		OnError("Couldn't set speaking", err)
	}

	// Stop speaking when done (voice overlay feature)
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			OnError("Couldn't stop speaking", err)
		}
		// Make sure processes are cleaned up
		ytDlpCmd.Process.Kill()
		ffmpegCmd.Process.Kill()
	}()

	send := make(chan []int16, 2)
	defer close(send)

	closeCh := make(chan bool)
	go func() {
		SendPCM(v, send)
		closeCh <- true
	}()

	// Add a minimum playback timer for very short clips
	minPlayTimer := time.NewTimer(500 * time.Millisecond)
	defer minPlayTimer.Stop()

	dataReceived := false

	// Track pause state
	isPaused := false

	// Stream audio from ffmpeg
	for {
		// Check pause channel
		select {
		case newState := <-pauseCh:
			isPaused = newState
			continue
		default:
			// No new pause state, continue
		}

		// If paused, wait and check again
		if isPaused {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Process audio normally
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if !dataReceived {
				// If we never got any data, wait a bit more
				select {
				case <-minPlayTimer.C:
					return
				case <-closeCh:
					return
				}
			}
			return
		}
		if err != nil {
			OnError("Error reading from ffmpeg stdout", err)
			return
		}

		dataReceived = true

		// Apply volume adjustment
		// Use v.GuildID for per-guild volume
		volumeMutex.Lock()
		currentVolume, ok := volume[v.GuildID]
		if !ok {
			currentVolume = 1.0
			volume[v.GuildID] = 1.0
		}
		volumeMutex.Unlock()

		for i := range audiobuf {
			// Calculate new value and clamp to int16 range to prevent distortion
			newValue := float64(audiobuf[i]) * currentVolume
			if newValue > maxClampValue {
				newValue = maxClampValue
			} else if newValue < minClampValue {
				newValue = minClampValue
			}
			audiobuf[i] = int16(newValue)
		}

		// Send audio data to channel
		select {
		case send <- audiobuf:
		case <-closeCh:
			return
		}
	}
}
