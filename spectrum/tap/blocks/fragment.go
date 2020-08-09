package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/storage"
)

// Fragment block for storing data without flag or checksum byte.
// Case #6: fragment data blocks.
//
// Fragmented data blocks cannot be produced by ROM saving routines - they must
// be produced with machine code programs. They have less than 2 bytes. In some
// games you really find zero length fragment data blocks.
type Fragment struct {
	Length uint16 // Length of the data in this block: 0 or 1 byte.
	Data   []byte // The essential data (may be empty)
}

// Read block data - reads 1 byte unless fragment size is zero length.
// It is expected that the tape pointer is at the correct position for reading.
func (b *Fragment) Read(reader *storage.Reader) {
	b.Length = reader.ReadShort()
	if b.Length > 0 {
		b.Data = reader.ReadBytes(int(b.Length))
	}
}

func (b Fragment) Id() uint8 {
	return 0x6 // FIXME: is this the correct ID to return?
}

func (b Fragment) Name() string {
	return "Data Fragment"
}

func (b Fragment) Filename() string {
	return ""
}

func (b Fragment) BlockData() []byte {
	return b.Data
}

// String returns a formatted string for the block
func (b Fragment) String() string {
	return fmt.Sprintf("%-13s: %d bytes", b.Name(), len(b.Data))
}
