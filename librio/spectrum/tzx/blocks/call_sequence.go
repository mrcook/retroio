package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// CallSequence
// ID: 26h (38d)
// This block is an analogue of the CALL Subroutine statement. It basically executes a sequence of
// blocks that are somewhere else and then goes back to the next block. Because more than one call
// can be normally used you can include a list of sequences to be called. The 'nesting' of call
// blocks is also not allowed for the simplicity reasons. You can, of course, use the CALL blocks
// in the LOOP sequences and vice versa. The value is relative for the obvious reasons - so that
// you can add some blocks in the beginning of the file without disturbing the call values. Please
// take a look at 'Jump To Block' for reference on the values.
type CallSequence struct {
	BlockID types.BlockType
	Count   uint16   // Number of calls to be made
	Blocks  []uint16 // Array of call block numbers (relative-signed offsets)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CallSequence) Read(reader *storage.Reader) error {
	c.BlockID = types.BlockType(reader.ReadByte())
	if c.BlockID != c.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", c.Id(), c.BlockID)
	}

	c.Count = reader.ReadShort()

	for i := 0; i < int(c.Count); i++ {
		c.Blocks = append(c.Blocks, reader.ReadShort())
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (c CallSequence) Id() types.BlockType {
	return types.CallSequence
}

// Name of the block as given in the TZX specification.
func (c CallSequence) Name() string {
	return "Call Sequence"
}
func (c CallSequence) BlockData() tap.Block {
	return nil
}
