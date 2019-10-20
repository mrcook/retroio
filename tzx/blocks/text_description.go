package blocks

import (
	"bufio"
	"fmt"

	"mrcook/tzxit/tape"
)

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

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (t *TextDescription) Read(reader *bufio.Reader) {
	t.Length, _ = reader.ReadByte()

	for _, b := range tape.ReadNextBytes(reader, int(t.Length)) {
		t.Description = append(t.Description, b)
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (t TextDescription) Id() uint8 {
	return 0x30
}

// Name of the block as given in the TZX specification.
func (t TextDescription) Name() string {
	return "Text Description"
}

// ToString returns a human readable string of the block data
func (t TextDescription) ToString() string {
	return fmt.Sprintf("> %-19s : %s", t.Name(), t.Description)
}
