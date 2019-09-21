package blocks

import (
	"bufio"
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// SequenceOfPulses
// ID: 13h (19d)
// This will produce N pulses, each having its own timing. Up to 255 pulses can be stored in this
// block; this is useful for non-standard sync tones used by some protection schemes.
type SequenceOfPulses struct {
	Count   uint8    // N BYTE  Number of pulses
	Lengths []uint16 // WORD[N] Pulses' lengths
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *SequenceOfPulses) Read(reader *bufio.Reader) {
	s.Count, _ = reader.ReadByte()

	for i := 0; i < int(s.Count); i++ {
		s.Lengths = append(s.Lengths, tape.ReadShort(reader))
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s SequenceOfPulses) Id() uint8 {
	return 0x13
}

// Name of the block as given in the TZX specification.
func (s SequenceOfPulses) Name() string {
	return "Sequence of Pulses"
}

// ToString returns a human readable string of the block data
func (s SequenceOfPulses) ToString() string {
	return fmt.Sprintf("> %-19s : %d pulses", s.Name(), s.Count)
}
