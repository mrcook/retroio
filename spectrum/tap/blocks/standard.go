package blocks

import (
	"fmt"
	"io"

	"retroio/storage"
)

// Standard data block for storing (2+[data length]) bytes.
// Case #5: standard data blocks or custom data blocks.
type Standard struct {
	Length uint16 // Length of the data in this block

	Flag     uint8  // Always 255 indicating a standard ROM loading data block or any other value to build a custom data block
	Data     []byte // The essential data (may be empty)
	Checksum uint8  // Simply all bytes (including flag byte) XORed
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *Standard) Read(reader *storage.Reader) {
	b.Length = reader.ReadShort()
	b.Flag = reader.ReadByte()

	b.Data = make([]byte, b.Length-2)
	_, err := reader.Read(b.Data)
	if err != nil && err != io.EOF {
		return
	}

	b.Checksum = reader.ReadByte()
}

func (b Standard) Id() uint8 {
	return b.Flag
}

func (b Standard) Name() string {
	return "Standard Data"
}

func (b Standard) Filename() string {
	return ""
}

func (b Standard) BlockData() []byte {
	return b.Data
}

// String returns a formatted string for the block
func (b Standard) String() string {
	return fmt.Sprintf("%-13s: %d bytes", b.Name(), len(b.Data))
}
