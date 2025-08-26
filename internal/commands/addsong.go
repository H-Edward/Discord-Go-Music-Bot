package commands

import (
	"discord-go-music-bot/internal/audio"
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddSong(s *discordgo.Session, m *discordgo.MessageCreate, search_mode bool) { // mode (false for play, true for search)
	var url string

	if !discordutil.IsUserInVoiceChannel(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You must be in a voice channel to use this command.")
		return
	}

	if search_mode {
		var hadToSanitise bool
		if len(m.Content) < 7 {
			s.ChannelMessageSend(m.ChannelID, "Invalid search query")
			return
		}

		searchQuery := strings.TrimSpace(m.Content[8:])

		if !validation.IsValidSearchQuery(searchQuery) {
			var searchQuerySafeToUse bool
			searchQuery, searchQuerySafeToUse = validation.SanitiseSearchQuery(searchQuery)
			hadToSanitise = true
			if !searchQuerySafeToUse {
				s.ChannelMessageSend(m.ChannelID, "Invalid search query")
				return
			}
		}

		var found_result bool
		url, found_result = audio.SearchYoutube(searchQuery)

		if !found_result {
			log.Println(constants.ANSIRed + "No results found for: " + searchQuery + constants.ANSIReset)
			s.ChannelMessageSend(m.ChannelID, "No results found for: "+searchQuery)
			return
		}

		if hadToSanitise {
			s.ChannelMessageSend(m.ChannelID, "Found: "+url+" using: "+searchQuery)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Found: "+url)
		}
	} else {
		if len(m.Content) < 6 {
			s.ChannelMessageSend(m.ChannelID, "Invalid URL")
			return
		}

		url = strings.TrimSpace(m.Content[6:])

		if !validation.IsValidURL(url) {
			s.ChannelMessageSend(m.ChannelID, "Invalid URL")
			return
		}

	}
	state.QueueMutex.Lock()
	state.Queue[m.GuildID] = append(state.Queue[m.GuildID], url)
	state.QueueMutex.Unlock()

	state.PlayingMutex.Lock()
	isAlreadyPlaying := state.Playing[m.GuildID]
	state.PlayingMutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, "Added to queue.")

	if !isAlreadyPlaying {
		// Start processing the queue if the bot is idle
		state.PlayingMutex.Lock()
		state.Playing[m.GuildID] = true
		state.PlayingMutex.Unlock()
		audio.ProcessQueue(s, m.GuildID, m.ChannelID, m)
	}
}
