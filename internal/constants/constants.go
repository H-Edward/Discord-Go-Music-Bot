package constants

const (
	// magic numbers
	Channels  int = 2                   // 1 for mono, 2 for stereo
	FrameRate int = 48000               // audio sampling rate
	FrameSize int = 960                 // uint16 size of each audio frame
	MaxBytes  int = (FrameSize * 2) * 2 // max size of opus data

	FfmpegBufferSize int = 16384 // 2^14

	MaxClampValue float64 = 32767  // max volume
	MinClampValue float64 = -32768 // min volume

	// ANSI color codes for terminal output
	ANSIBold   = "\033[1m"
	ANSIBlue   = "\033[34m"
	ANSICyan   = "\033[36m"
	ANSIYellow = "\033[33m"
	ANSIGreen  = "\033[32m"
	ANSIRed    = "\033[31m"
	ANSIReset  = "\033[0m"
)
