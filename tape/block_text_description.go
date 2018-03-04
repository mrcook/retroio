package tape

import "fmt"

// TextDescription
// ID: 30h (48d)
// This is meant to identify parts of the tape, so you know where level 1 starts, where to rewind
// to when the game ends, etc. This description is not guaranteed to be shown while the tape is
// playing, but can be read while browsing the tape or changing the tape pointer.
// The description can be up to 255 characters long but please keep it down to about 30 so the
// programs can show it in one line (where this is appropriate).
type TextDescription struct {
	Length      uint8  // N BYTE  Length of the text description
	Description []byte // CHAR[N] Text description in ASCII format
}

func (t *TextDescription) Process(file *File) {
	t.Length, _ = file.ReadByte()

	for _, b := range file.ReadBytes(int(t.Length)) {
		t.Description = append(t.Description, b)
	}
}

func (t TextDescription) Id() int {
	return 48
}

func (t TextDescription) Name() string {
	return "Text Description"
}

// Metadata returns a human readable string of the block data
func (t TextDescription) Metadata() string {
	str := ""
	str += fmt.Sprintf("Length:      %d\n", t.Length)
	str += fmt.Sprintf("Description: %s\n", t.Description)
	return str
}
