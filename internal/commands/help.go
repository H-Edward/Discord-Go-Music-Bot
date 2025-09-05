package commands

import (
	"discord-go-music-bot/internal/state"
)

func Help(ctx state.Context) {
	helpMessage := "Commands:\n" +
		"/ping - Responds with Pong\n" +
		"/pong - Responds with Ping\n" +
		"/play <url> - Plays a song from the given URL\n" +
		"/search <query> - Searches for a song and plays it\n" +
		"/skip - Skips the current song\n" +
		"/queue - Shows the current queue\n" +
		"/stop - Stops playback and clears the queue\n" +
		"/pause - Pauses playback\n" +
		"/resume - Resumes playback\n" +
		"/volume <value> - Sets the volume (0 to 200)\n" +
		"/currentvolume - Shows the current volume\n" +
		"/nuke <number> - Deletes the specified number of messages\n" +
		"/uptime - Shows how long the bot has been running\n" +
		"/version - Shows a hash-based version of the bot\n" +
		"/oss - Provides a link to the bot's source code\n" +
		"/help - Shows this help message\n"
	ctx.Reply(helpMessage)
}
