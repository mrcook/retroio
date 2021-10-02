package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// GroupStart
// ID: 21h (33d)
// This block marks the start of a group of blocks which are to be treated as one single
// (composite) block. This is very handy for tapes that use lots of subblocks like Bleepload
// (which may well have over 160 custom loading blocks). You can also give the group a name
// (example 'Bleepload Block 1').
// For each group start block, there must be a group end block. Nesting of groups is not allowed.
type GroupStart struct {
	BlockID   types.BlockType
	Length    uint8  // Length of the group name string
	GroupName []byte // Group name in ASCII format (please keep it under 30 characters long)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (g *GroupStart) Read(reader *storage.Reader) error {
	g.BlockID = types.BlockType(reader.ReadByte())
	if g.BlockID != g.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", g.Id(), g.BlockID)
	}

	g.Length = reader.ReadByte()

	for _, b := range reader.ReadBytes(int(g.Length)) {
		g.GroupName = append(g.GroupName, b)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (g GroupStart) Id() types.BlockType {
	return types.GroupStart
}

// Name of the block as given in the TZX specification.
func (g GroupStart) Name() string {
	return "Group Start"
}

func (g GroupStart) BlockData() tap.Block {
	return nil
}
