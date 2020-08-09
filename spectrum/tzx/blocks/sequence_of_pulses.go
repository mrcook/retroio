package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// SequenceOfPulses
// ID: 13h (19d)
// This will produce N pulses, each having its own timing. Up to 255 pulses can be stored in this
// block; this is useful for non-standard sync tones used by some protection schemes.
type SequenceOfPulses struct {
	BlockID types.BlockType
	Count   uint8    // Number of pulses
	Lengths []uint16 // Pulses' lengths
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *SequenceOfPulses) Read(reader *storage.Reader) error {
	s.BlockID = types.BlockType(reader.ReadByte())
	if s.BlockID != s.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", s.Id(), s.BlockID)
	}

	s.Count = reader.ReadByte()

	for i := 0; i < int(s.Count); i++ {
		s.Lengths = append(s.Lengths, reader.ReadShort())
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s SequenceOfPulses) Id() types.BlockType {
	return types.SequenceOfPulses
}

// Name of the block as given in the TZX specification.
func (s SequenceOfPulses) Name() string {
	return "Sequence of Pulses"
}

func (s SequenceOfPulses) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (s SequenceOfPulses) String() string {
	return fmt.Sprintf("%-19s : %d pulses", s.Name(), s.Count)
}
