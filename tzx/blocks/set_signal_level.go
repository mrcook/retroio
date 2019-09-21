package blocks

import (
	"bufio"
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// SetSignalLevel
// ID: 2Bh (43d)
// This block sets the current signal level to the specified value (high or low). It should be used
// whenever it is necessary to avoid any ambiguities, e.g. with custom loaders which are level-sensitive.
type SetSignalLevel struct {
	Length      uint32 // DWORD Block length (without these four bytes)
	SignalLevel uint8  // BYTE  Signal level (0=low, 1=high)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *SetSignalLevel) Read(reader *bufio.Reader) {
	s.Length = tape.ReadLong(reader)
	s.SignalLevel, _ = reader.ReadByte()
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s SetSignalLevel) Id() uint8 {
	return 0x2b
}

// Name of the block as given in the TZX specification.
func (s SetSignalLevel) Name() string {
	return "Set Signal Level"
}

// ToString returns a human readable string of the block data
func (s SetSignalLevel) ToString() string {
	return fmt.Sprintf("> %-19s : signal level: %d", s.Name(), s.SignalLevel)
}
