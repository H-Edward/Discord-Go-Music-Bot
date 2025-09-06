package commands

import (
	"discord-go-music-bot/internal/state"
	"fmt"
	"strings"
)

func ShowQueue(ctx state.Context) {
	state.QueueMutex.Lock()
	defer state.QueueMutex.Unlock()

	if len(state.Queue[ctx.GetGuildID()]) == 0 {
		ctx.Reply("Queue is empty.")
		return
	}

	// Make a formatted list of songs, "[N] URL""
	var formattedQueue []string
	for i, song := range state.Queue[ctx.GetGuildID()] {
		formattedQueue = append(formattedQueue, fmt.Sprintf("[%d] %s", i+1, song))
	}

	ctx.Reply("Current queue:\n" + strings.Join(formattedQueue, "\n"))
}
