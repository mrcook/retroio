package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// Select
// ID: 28h (40d)
// This block is useful when the tape consists of two or more separately-loadable parts. With this
// block, you are able to select one of the parts and the utility/emulator will start loading from
// that block. For example you can use it when the game has a separate Trainer or when it is a
// multi-load. Of course, to make some use of it the emulator/utility has to show a menu with the
// selections when it encounters such a block. All offsets are relative signed words.
type Select struct {
	Length     uint16      // WORD  Length of the whole block (without these two bytes)
	Count      uint8       // N BYTE  Number of selections
	Selections []Selection // SELECT[N] List of selections
}

type Selection struct {
	RelativeOffset int16   // WORD  Relative Offset as `signed` value
	Length         uint8   // L BYTE  Length of description text
	Description    []uint8 // CHAR[L] Description text (please use single line and max. 30 chars)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *Select) Read(file *tape.File) {
	s.Length = file.ReadShort()
	s.Count, _ = file.ReadByte()

	for i := 0; i < int(s.Count); i++ {
		var selection Selection
		selection.RelativeOffset = file.ReadSignedShort()
		selection.Length, _ = file.ReadByte()
		for _, b := range file.ReadBytes(int(selection.Length)) {
			selection.Description = append(selection.Description, b)
		}
		s.Selections = append(s.Selections, selection)
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s Select) Id() uint8 {
	return 0x28
}

// Name of the block as given in the TZX specification.
func (s Select) Name() string {
	return "Select"
}

// ToString returns a human readable string of the block data
func (s Select) ToString() string {
	str := fmt.Sprintf("> %s\n", s.Name())
	for _, b := range s.Selections {
		str += fmt.Sprintf("- Offset:      %d\n", b.RelativeOffset)
		str += fmt.Sprintf("  Description: %s\n", b.Description)
	}
	return str
}
