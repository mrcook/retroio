package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// PureTone
// ID: 12h (18d)
// This will produce a tone which is basically the same as the pilot tone in the ID 10, ID 11
// blocks. You can define how long the pulse is and how many pulses are in the tone.
type PureTone struct {
	BlockID    types.BlockType
	Length     uint16 // Length of one pulse in T-states
	PulseCount uint16 // Number of pulses
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PureTone) Read(reader *storage.Reader) error {
	p.BlockID = types.BlockType(reader.ReadByte())
	if p.BlockID != p.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", p.Id(), p.BlockID)
	}

	p.Length = reader.ReadShort()
	p.PulseCount = reader.ReadShort()

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PureTone) Id() types.BlockType {
	return types.PureTone
}

// Name of the block as given in the TZX specification.
func (p PureTone) Name() string {
	return "Pure Tone"
}

func (p PureTone) BlockData() tap.Block {
	return nil
}
