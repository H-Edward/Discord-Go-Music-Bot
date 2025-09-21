package audio

import (
	"bufio"
	"encoding/binary"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"
)

// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayURL(v *discordgo.VoiceConnection, url string, stop <-chan bool, pauseCh <-chan bool) {

	if !validation.IsValidURL(url) {
		OnError("Invalid URL"+url, nil)
		return
	}

	// Create the yt-dlp command to download only the audio with best quality
	ytDlpCmd := exec.Command("yt-dlp",
		"-f", "bestaudio",
		"--no-playlist",
		"-o", "-",
		url) // Get only audio, best quality

	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(constants.FrameRate), "-ac", strconv.Itoa(constants.Channels), "pipe:1")

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
	ffmpegbuf := bufio.NewReaderSize(ffmpegOut, constants.FfmpegBufferSize)

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
		case <-time.After(3 * time.Hour): // Fallback timeout
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
		audiobuf := make([]int16, constants.FrameSize*constants.Channels)


		// Process audio normally
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
		state.VolumeMutex.Lock()
		currentVolume, ok := state.Volume[v.GuildID]
		if !ok {
			currentVolume = 1.0
			state.Volume[v.GuildID] = 1.0
		}
		state.VolumeMutex.Unlock()

		for i := range audiobuf {
			// Calculate new value and clamp to int16 range to prevent distortion
			newValue := float64(audiobuf[i]) * currentVolume
			if newValue > constants.MaxClampValue {
				newValue = constants.MaxClampValue
			} else if newValue < constants.MinClampValue {
				newValue = constants.MinClampValue
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
