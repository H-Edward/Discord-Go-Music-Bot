package state

import (
	"sync"
	"time"

	"layeh.com/gopus"
)

func init() {
	StartTime = time.Now()
}

var (
	Token        string
	Queue        = make(map[string][]string) // Guild ID -> Queue of URLs
	QueueMutex   sync.Mutex
	Playing      = make(map[string]bool)
	PlayingMutex sync.Mutex
	Paused       = make(map[string]bool) // Guild ID -> Pause state
	PauseMutex   sync.Mutex
	Volume       = make(map[string]float64) // Guild ID -> Volume
	VolumeMutex  sync.Mutex
	OpusEncoder  *gopus.Encoder
	StopChannels = make(map[string]chan bool)
	StopMutex    sync.Mutex
	PauseChs     = make(map[string]chan bool) // Map of guild ID to pause channels
	PauseChMutex sync.Mutex

	StartTime time.Time // Time when the bot started

	GoSourceHash string // short hash of all go source files for !version
)
