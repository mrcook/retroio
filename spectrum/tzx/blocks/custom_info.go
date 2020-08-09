package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// CustomInfo
// ID: 35h (53d)
// This block can be used to save any information you want. For example, it might contain some
// information written by a utility, extra settings required by a particular emulator, or even
// poke data.
type CustomInfo struct {
	BlockID        types.BlockType
	Identification [10]byte // Identification string (in ASCII)
	Length         uint32   // Length of the custom info
	Info           []uint8  // Custom info
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CustomInfo) Read(reader *storage.Reader) error {
	c.BlockID = types.BlockType(reader.ReadByte())
	if c.BlockID != c.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", c.Id(), c.BlockID)
	}

	for i, b := range reader.ReadBytes(10) {
		c.Identification[i] = b
	}

	c.Length = reader.ReadLong()

	for _, b := range reader.ReadBytes(int(c.Length)) {
		c.Info = append(c.Info, b)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (c CustomInfo) Id() types.BlockType {
	return types.CustomInfo
}

// Name of the block as given in the TZX specification.
func (c CustomInfo) Name() string {
	return "Custom Info"
}

func (c CustomInfo) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (c CustomInfo) String() string {
	return fmt.Sprintf("%-19s : %s - %s", c.Name(), c.Identification, c.Info)
}
