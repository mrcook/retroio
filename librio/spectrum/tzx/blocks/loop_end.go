package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// LoopEnd
// ID: 25h (37d)
// This is the same as BASIC's NEXT statement. It means that the utility should jump back to the
// start of the loop if it hasn't been run for the specified number of times.
// This block has no body.
type LoopEnd struct {
	BlockID types.BlockType
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (l *LoopEnd) Read(reader *storage.Reader) error {
	l.BlockID = types.BlockType(reader.ReadByte())
	if l.BlockID != l.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", l.Id(), l.BlockID)
	}
	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (l LoopEnd) Id() types.BlockType {
	return types.LoopEnd
}

// Name of the block as given in the TZX specification.
func (l LoopEnd) Name() string {
	return "Loop End"
}

func (l LoopEnd) BlockData() tap.Block {
	return nil
}
