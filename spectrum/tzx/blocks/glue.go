package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// GlueBlock
// ID: 5Ah (90d)
// This block is generated when you merge two ZX Tape files together. It is here so that you can
// easily copy the files together and use them. Of course, this means that resulting file would
// be 10 bytes longer than if this block was not used. All you have to do if you encounter this
// block ID is to skip next 9 bytes.
// If you can avoid using this block for this purpose, then do so; it is preferable to use a
// utility to join the two files and ensure that they are both of the higher version number.
type GlueBlock struct {
	BlockID types.BlockType
	Value   [9]byte // Value: { "XTape!",0x1A,MajR,MinR } Just skip these 9 bytes and you will end up on the next ID.
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (g *GlueBlock) Read(reader *storage.Reader) error {
	g.BlockID = types.BlockType(reader.ReadByte())
	if g.BlockID != g.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", g.Id(), g.BlockID)
	}

	for i, b := range reader.ReadBytes(9) {
		g.Value[i] = b
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (g GlueBlock) Id() types.BlockType {
	return types.GlueBlock
}

// Name of the block as given in the TZX specification.
func (g GlueBlock) Name() string {
	return "Glue Block"
}

func (g GlueBlock) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (g GlueBlock) String() string {
	return fmt.Sprintf("%s", g.Name())
}
