package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
)

// TextDescription
// ID: 30h (48d)
// This is meant to identify parts of the tape, so you know where level 1 starts, where to rewind
// to when the game ends, etc. This description is not guaranteed to be shown while the tape is
// playing, but can be read while browsing the tape or changing the tape pointer.
// The description can be up to 255 characters long but please keep it down to about 30 so the
// programs can show it in one line (where this is appropriate).
type TextDescription struct {
	BlockID     types.BlockType
	Length      uint8  // Length of the text description
	Description []byte // Text description in ASCII format
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (t *TextDescription) Read(reader *storage.Reader) error {
	t.BlockID = types.BlockType(reader.ReadByte())
	if t.BlockID != t.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", t.Id(), t.BlockID)
	}

	t.Length = reader.ReadByte()

	for _, b := range reader.ReadBytes(int(t.Length)) {
		t.Description = append(t.Description, b)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (t TextDescription) Id() types.BlockType {
	return types.TextDescription
}

// Name of the block as given in the TZX specification.
func (t TextDescription) Name() string {
	return "Text Description"
}

func (t TextDescription) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (t TextDescription) String() string {
	return fmt.Sprintf("%-19s : %s", t.Name(), t.Description)
}
