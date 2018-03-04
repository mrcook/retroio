package tape

import "fmt"

// SequenceOfPulses
// ID: 13h (19d)
// This will produce N pulses, each having its own timing. Up to 255 pulses can be stored in this
// block; this is useful for non-standard sync tones used by some protection schemes.
type SequenceOfPulses struct {
	Count   uint8    // N BYTE  Number of pulses
	Lengths []uint16 // WORD[N] Pulses' lengths
}

func (s *SequenceOfPulses) Process(file *File) {
	s.Count, _ = file.ReadByte()

	for i := 0; i < int(s.Count); i++ {
		s.Lengths = append(s.Lengths, file.ReadShort())
	}
}

func (s SequenceOfPulses) Id() int {
	return 19
}

func (s SequenceOfPulses) Name() string {
	return "Sequence of Pulses"
}

// Metadata returns a human readable string of the block data
func (s SequenceOfPulses) Metadata() string {
	return fmt.Sprintf("> %-19s : %d pulses", s.Name(), s.Count)
}
