package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// GroupEnd
// ID: 22h (34d)
// This indicates the end of a group. This block has no body.
type GroupEnd struct {
	BlockID types.BlockType
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (g *GroupEnd) Read(reader *storage.Reader) error {
	g.BlockID = types.BlockType(reader.ReadByte())
	if g.BlockID != g.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", g.Id(), g.BlockID)
	}
	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (g GroupEnd) Id() types.BlockType {
	return types.GroupEnd
}

// Name of the block as given in the TZX specification.
func (g GroupEnd) Name() string {
	return "Group End"
}

func (g GroupEnd) BlockData() tap.Block {
	return nil
}
