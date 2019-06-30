package tap

import "fmt"

// AlphanumericHeader for storing string arrays
// case #3: alphanumeric data array header
type AlphanumericHeader struct {
	Flag         uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // BYTE      Always 2: Byte indicating an alphanumeric array.
	Filename     [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // WORD      Length of data following the header = length of string array + 3.
	UnusedByte   uint8    // BYTE      Unused byte.
	VariableName byte     // BYTE      = (1..26 meaning A$..Z$) + 192.
	UnusedWord   uint16   // WORD      = 32768.
	Checksum     uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

// ToString returns a formatted string for the header
func (h AlphanumericHeader) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", flagLabel(h.Flag), typeLabel(h.DataType))
	str += fmt.Sprintf("     - Filename     : %s\n", h.Filename)
	str += fmt.Sprintf("     - Variable Name: %c", h.VariableName-192)
	return str
}
