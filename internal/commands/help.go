package commands

import (
	"discord-go-music-bot/internal/state"
)

func Help(ctx *state.Context) {
	prefix := "/"
	if ctx.GetSourceType() == int(state.SourceTypeMessage) {
		prefix = "!"
	}

	helpMessage := "Commands:\n" +
		prefix + "ping - Responds with Pong\n" +
		prefix + "pong - Responds with Ping\n" +
		prefix + "play <url> - Plays a song from the given URL\n" +
		prefix + "search <query> - Searches for a song and plays it\n" +
		prefix + "skip - Skips the current song\n" +
		prefix + "queue - Shows the current queue\n" +
		prefix + "stop - Stops playback and clears the queue\n" +
		prefix + "pause - Pauses playback\n" +
		prefix + "resume - Resumes playback\n" +
		prefix + "volume <value> - Sets the volume (0 to 200)\n" +
		prefix + "currentvolume - Shows the current volume\n" +
		prefix + "nuke <number> - Deletes the specified number of messages\n" +
		prefix + "uptime - Shows how long the bot has been running\n" +
		prefix + "version - Shows a hash-based version of the bot\n" +
		prefix + "oss - Provides a link to the bot's source code\n" +
		prefix + "help - Shows this help message\n"
	ctx.Reply(helpMessage)
}
