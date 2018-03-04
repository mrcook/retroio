package tape

import "fmt"

// StopTapeWhen48kMode
// ID: 2Ah (42d)
// When this block is encountered, the tape will stop ONLY if the machine is an 48K Spectrum.
// This block is to be used for multi-loading games that load one level at a time in 48K mode,
// but load the entire tape at once if in 128K mode.
// This block has no body of its own, but follows the extension rule.
type StopTapeWhen48kMode struct {
	Length uint32 // DWORD Length of the block without these four bytes (0)
}

func (s *StopTapeWhen48kMode) Process(file *File) {
	s.Length = file.ReadLong()
}

func (s StopTapeWhen48kMode) Id() int {
	return 42
}

func (s StopTapeWhen48kMode) Name() string {
	return "Stop Tape when in 48k Mode"
}

// Metadata returns a human readable string of the block data
func (s StopTapeWhen48kMode) Metadata() string {
	return fmt.Sprintf("Length: %d\n", s.Length)
}
