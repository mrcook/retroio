package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
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

func (g GroupStart) BlockData() tap.BlockI {
	return nil
}

// String returns a human readable string of the block data
func (g GroupStart) String() string {
	return fmt.Sprintf("%-19s : %s", g.Name(), g.GroupName)
}

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

func (g GroupEnd) BlockData() tap.BlockI {
	return nil
}

// String returns a human readable string of the block data
func (g GroupEnd) String() string {
	return fmt.Sprintf("%s", g.Name())
}
