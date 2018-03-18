package tap

import "fmt"

// StandardDataBlocks for storing (2+[data length]) bytes
// case #5: standard data blocks or custom data blocks
type StandardDataBlocks struct {
	Flag     uint8  // BYTE    flag byte      always 255 indicating a standard ROM loading data block or any other value to build a custom data block
	Data     []byte // BYTE[]  [data length]  data block the essential data (may be empty)
	Checksum uint8  // BYTE    checksum byte  simply all bytes (including flag byte) XORed
}

// ToString returns a formatted string for the block
func (b StandardDataBlocks) ToString() string {
	return fmt.Sprintf("  - %s\n", flagLabel(b.Flag))
}
