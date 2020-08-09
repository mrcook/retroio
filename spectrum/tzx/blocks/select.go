package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// Select
// ID: 28h (40d)
// This block is useful when the tape consists of two or more separately-loadable parts. With this
// block, you are able to select one of the parts and the utility/emulator will start loading from
// that block. For example you can use it when the game has a separate Trainer or when it is a
// multi-load. Of course, to make some use of it the emulator/utility has to show a menu with the
// selections when it encounters such a block. All offsets are relative signed words.
type Select struct {
	BlockID    types.BlockType
	Length     uint16      // Length of the whole block (without these two bytes)
	Count      uint8       // Number of selections
	Selections []Selection // List of selections
}

type Selection struct {
	RelativeOffset int16   // Relative Offset as `signed` value
	Length         uint8   // Length of description text
	Description    []uint8 // Description text (please use single line and max. 30 chars)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *Select) Read(reader *storage.Reader) error {
	s.BlockID = types.BlockType(reader.ReadByte())
	if s.BlockID != s.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", s.Id(), s.BlockID)
	}

	s.Length = reader.ReadShort()
	s.Count = reader.ReadByte()

	for i := 0; i < int(s.Count); i++ {
		var selection Selection
		selection.RelativeOffset = int16(reader.ReadShort())
		selection.Length = reader.ReadByte()
		for _, b := range reader.ReadBytes(int(selection.Length)) {
			selection.Description = append(selection.Description, b)
		}
		s.Selections = append(s.Selections, selection)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s Select) Id() types.BlockType {
	return types.Select
}

// Name of the block as given in the TZX specification.
func (s Select) Name() string {
	return "Select"
}

func (s Select) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (s Select) String() string {
	str := fmt.Sprintf("%s\n", s.Name())
	for _, b := range s.Selections {
		str += fmt.Sprintf("- Offset:      %d\n", b.RelativeOffset)
		str += fmt.Sprintf("  Description: %s\n", b.Description)
	}
	return str
}
