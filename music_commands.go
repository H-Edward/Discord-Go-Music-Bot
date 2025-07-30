package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func showQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	if len(queue[m.GuildID]) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Queue is empty.")
		return
	}
	// Make a formatted list of songs, "[N] URL""
	queueMutex.Lock()
	var formattedQueue []string
	for i, song := range queue[m.GuildID] {
		formattedQueue = append(formattedQueue, fmt.Sprintf("[%d] %s", i+1, song))
	}
	queueMutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, "Current queue:\n"+strings.Join(formattedQueue, "\n"))
}

func pauseSong(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID

	// Check if the bot is in a voice channel
	if !botInChannel(s, guildID) {
		s.ChannelMessageSend(m.ChannelID, "Not in a voice channel.")
		return
	}

	pauseMutex.Lock()
	currentState := paused[guildID]
	paused[guildID] = !currentState // Toggle pause state
	pauseMutex.Unlock()

	// Signal the pause channel with the new state
	pauseChMutex.Lock()
	if ch, exists := pauseChs[guildID]; exists {
		select {
		case ch <- !currentState: // Send the new state
		default: // Channel is full, discard
			// This prevents blocking if the channel is full
		}
	}
	pauseChMutex.Unlock()

	if currentState {
		s.ChannelMessageSend(m.ChannelID, "Resumed playback.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Paused playback.")
	}
}

func setVolume(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !volume <value between 0 and 200>")
		return
	}

	newVolume, err := strconv.ParseFloat(args[1], 64)
	if err != nil || newVolume < 0.0 || newVolume > 200.0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid volume value. Please specify a number between 0 and 200.")
		return
	}
	var preservedVolume float64 = newVolume
	// Normalize the volume to a range of 0.0 to 2.0
	newVolume = newVolume / 100.0 // Convert percentage to a factor


	volumeMutex.Lock()
	if _, ok := volume[m.GuildID]; !ok {
		volume[m.GuildID] = 1.0 // Initialize to default if not set
	}
	volume[m.GuildID] = newVolume
	volumeMutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Volume set to %.2f%%", preservedVolume))
}

func skipSong(s *discordgo.Session, m *discordgo.MessageCreate) {
	vc, err := getVoiceConnection(s, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not in a voice channel")
		return
	}

	// Signal the current song to stop
	stopMutex.Lock()
	if stopChan, exists := stopChannels[m.GuildID]; exists {
		close(stopChan)
		delete(stopChannels, m.GuildID)
	}
	stopMutex.Unlock()

	vc.Speaking(false)

	s.ChannelMessageSend(m.ChannelID, "Skipping current song")

	// The song will stop, and the queue processor will automatically move to the next song
	// We don't need to start a new queue processor
}

func stopSong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get the voice connection for the guild
	vc, err := getVoiceConnection(s, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not in a voice channel")
		return
	}

	// Signal the current song to stop
	stopMutex.Lock()
	if stopChan, exists := stopChannels[m.GuildID]; exists {
		close(stopChan)
		delete(stopChannels, m.GuildID)
	}
	stopMutex.Unlock()

	// Clear the queue for the guild
	queueMutex.Lock()
	queue[m.GuildID] = []string{}
	queueMutex.Unlock()

	// Mark the bot as not playing
	playingMutex.Lock()
	playing[m.GuildID] = false
	playingMutex.Unlock()

	// Wait a moment for processes to terminate cleanly
	// then disconnect from the voice channel
	go func() {
		// Give a small delay for processes to clean up
		time.Sleep(500 * time.Millisecond)
		vc.Speaking(false)
		err = vc.Disconnect()
		if err != nil {
			log.Println(ANSIRed + "Error disconnecting from voice channel: " + err.Error() + ANSIReset)
		}
	}()

	// Notify the user
	s.ChannelMessageSend(m.ChannelID, "Stopped playback and cleared the queue.")
}

func currentVolume(s *discordgo.Session, m *discordgo.MessageCreate) {
	volumeMutex.Lock()
	defer volumeMutex.Unlock()

	currentVolume, ok := volume[m.GuildID]
	if !ok {
		currentVolume = 1.0 // Default volume if not set
		volume[m.GuildID] = 1.0
	}
	// Convert to percentage for display
	currentVolume = currentVolume * 100.0 // Convert factor back to percentage
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current volume is %.2f%%", currentVolume))
}
