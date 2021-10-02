package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// SetSignalLevel
// ID: 2Bh (43d)
// This block sets the current signal level to the specified value (high or low). It should be used
// whenever it is necessary to avoid any ambiguities, e.g. with custom loaders which are level-sensitive.
type SetSignalLevel struct {
	BlockID     types.BlockType
	Length      uint32 // Block length (without these four bytes)
	SignalLevel uint8  // Signal level (0=low, 1=high)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *SetSignalLevel) Read(reader *storage.Reader) error {
	s.BlockID = types.BlockType(reader.ReadByte())
	if s.BlockID != s.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", s.Id(), s.BlockID)
	}

	s.Length = reader.ReadLong()
	s.SignalLevel = reader.ReadByte()

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s SetSignalLevel) Id() types.BlockType {
	return types.SetSignalLevel
}

// Name of the block as given in the TZX specification.
func (s SetSignalLevel) Name() string {
	return "Set Signal Level"
}

func (s SetSignalLevel) BlockData() tap.Block {
	return nil
}
