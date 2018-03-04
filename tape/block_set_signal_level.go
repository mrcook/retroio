package tape

import "fmt"

// SetSignalLevel
// ID: 2Bh (43d)
// This block sets the current signal level to the specified value (high or low). It should be used
// whenever it is necessary to avoid any ambiguities, e.g. with custom loaders which are level-sensitive.
type SetSignalLevel struct {
	Length      uint32 // DWORD Block length (without these four bytes)
	SignalLevel uint8  // BYTE  Signal level (0=low, 1=high)
}

func (s *SetSignalLevel) Process(file *File) {
	s.Length = file.ReadLong()
	s.SignalLevel, _ = file.ReadByte()
}

func (s SetSignalLevel) Id() int {
	return 43
}

func (s SetSignalLevel) Name() string {
	return "Set Signal Level"
}

// Metadata returns a human readable string of the block data
func (s SetSignalLevel) Metadata() string {
	return fmt.Sprintf("Length: %d\n", s.Length)
}
