package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// ReturnFromSequence
// ID: 27h (39d)
// This block indicates the end of the Called Sequence. The next block played will be the block after
// the last CALL block (or the next Call, if the Call block had multiple calls).
// This block has no body.
type ReturnFromSequence struct {
	BlockID types.BlockType
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (r ReturnFromSequence) Read(reader *storage.Reader) error {
	r.BlockID = types.BlockType(reader.ReadByte())
	if r.BlockID != r.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", r.Id(), r.BlockID)
	}
	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (r ReturnFromSequence) Id() types.BlockType {
	return types.ReturnFromSequence
}

// Name of the block as given in the TZX specification.
func (r ReturnFromSequence) Name() string {
	return "Return from Sequence"
}

func (r ReturnFromSequence) BlockData() tap.Block {
	return nil
}
