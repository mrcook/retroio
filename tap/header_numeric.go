package tap

import "fmt"

// NumericHeader for storing numeric arrays
// case #2: numeric data array header
type NumericHeader struct {
	Flag         uint8    // BYTE     flag byte      always 0. Byte indicating a standard ROM loading header
	DataType     uint8    // BYTE     data type      always 1: Byte indicating a numeric array
	Filename     [10]byte // CHAR[10] file name      loading name of the program. filled with spaces (CHR$(32))
	DataLength   uint16   // WORD     [data length]  length of the following data (after the header) = length of number array * 5 +3
	UnusedByte   uint8    // BYTE     unused         unused byte
	VariableName byte     // BYTE     variable name  = (1..26 meaning A..Z) +128
	UnusedWord   uint16   // WORD     unused         = 32768
	Checksum     uint8    // BYTE     checksum byte  simply all bytes (including flag byte) XORed
}

// ToString returns a formatted string for the header
func (h NumericHeader) ToString() string {
	str := fmt.Sprintf("  - %-13s: %s\n", flagLabel(h.Flag), typeLabel(h.DataType))
	str += fmt.Sprintf("  - Filename     : %s\n", h.Filename)
	str += fmt.Sprintf("  - Variable Name: %c", h.VariableName-128)
	return str
}
