package commands

import (
	"discord-go-music-bot/internal/audio"
	"discord-go-music-bot/internal/constants"
	"discord-go-music-bot/internal/discordutil"
	"discord-go-music-bot/internal/logging"
	"discord-go-music-bot/internal/state"
	"discord-go-music-bot/internal/validation"
	"strings"
)

func AddSong(ctx *state.Context, search_mode bool) { // mode (false for play, true for search)
	var url string

	if !discordutil.IsUserInVoiceChannel(ctx) {
		ctx.Reply(constants.EmojiWarning + " You must be in a voice channel to use this command.")
		return
	}

	if search_mode {
		if ctx.SourceType == state.SourceTypeInteraction {
			// To avoid the discord timeout for interactions
			ctx.Reply(constants.EmojiSearch + " Searching...")
		}

		var hadToSanitise bool

		searchQuery := strings.TrimSpace(ctx.Arguments["query"])

		if !validation.IsValidSearchQuery(searchQuery) {
			var searchQuerySafeToUse bool
			searchQuery, searchQuerySafeToUse = validation.SanitiseSearchQuery(searchQuery)
			hadToSanitise = true
			if !searchQuerySafeToUse {
				ctx.Reply(constants.EmojiWarning + " Invalid search query")
				return
			}
		}

		var found_result bool
		url, found_result = audio.SearchYoutube(searchQuery)

		if !found_result {
			logging.ErrorLog("No results found for: " + searchQuery)
			ctx.Reply(constants.EmojiWarning + " No results found for: " + searchQuery)
			return
		}

		if hadToSanitise {
			ctx.Reply(constants.EmojiInfo + " Found: " + url + " using: " + searchQuery)
		} else {
			ctx.Reply(constants.EmojiInfo + " Found: " + url)
		}
	} else {
		if len(ctx.Arguments["url"]) < 6 {
			ctx.Reply(constants.EmojiWarning + " Invalid URL")
			return
		}

		url = strings.TrimSpace(ctx.Arguments["url"])

		if !validation.IsValidURL(url) {
			ctx.Reply(constants.EmojiWarning + " Invalid URL")
			return
		}

	}
	state.QueueMutex.Lock()
	state.Queue[ctx.GetGuildID()] = append(state.Queue[ctx.GetGuildID()], url)
	state.QueueMutex.Unlock()

	state.PlayingMutex.Lock()
	isAlreadyPlaying := state.Playing[ctx.GetGuildID()]
	state.PlayingMutex.Unlock()

	ctx.Reply(constants.EmojiSuccess + " Added to queue.")

	if !isAlreadyPlaying {
		// Start processing the queue if the bot is idle
		state.PlayingMutex.Lock()
		state.Playing[ctx.GetGuildID()] = true
		state.PlayingMutex.Unlock()
		audio.ProcessQueue(ctx)
	}
}
