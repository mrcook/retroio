package tape

import "fmt"

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

func (s *Select) Process(file *File) {
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

func (s Select) Id() int {
	return 40
}

func (s Select) Name() string {
	return "Select"
}

// Metadata returns a human readable string of the block data
func (s Select) Metadata() string {
	str := ""
	str += fmt.Sprintf("Length: %d\n", s.Length)
	str += fmt.Sprintf("Count:  %d\n", s.Count)

	str += fmt.Sprint("Selections:\n")
	for _, b := range s.Selections {
		str += fmt.Sprintf("- Offset:      %d\n", b.RelativeOffset)
		str += fmt.Sprintf("  Length:      %d\n", b.Length)
		str += fmt.Sprintf("  Description: %s\n", b.Description)
	}
	return str
}
