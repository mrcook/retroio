package blocks

import (
	"fmt"
	"io"

	"retroio/storage"
)

// Standard data block for storing (2+[data length]) bytes.
// Case #5: standard data blocks or custom data blocks.
type Standard struct {
	Length uint16 // WORD  Length of the data in this block

	Flag     uint8  // BYTE    flag byte      always 255 indicating a standard ROM loading data block or any other value to build a custom data block
	Data     []byte // BYTE[]  [data length]  data block the essential data (may be empty)
	Checksum uint8  // BYTE    checksum byte  simply all bytes (including flag byte) XORed
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *Standard) Read(reader *storage.Reader) {
	b.Length = reader.ReadShort()

	data := make([]byte, b.Length)
	_, _ = io.ReadFull(reader, data)

	b.Flag = data[0]
	b.Data = data[1 : len(data)-1]
	b.Checksum = data[len(data)-1]
}

func (b Standard) Id() uint8 {
	return b.Flag
}

func (b Standard) Name() string {
	return "Standard Data"
}

// ToString returns a formatted string for the block
func (b Standard) ToString() string {
	return fmt.Sprintf("%-13s: %d bytes", b.Name(), len(b.Data))
}
