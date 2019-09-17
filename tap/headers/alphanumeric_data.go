package headers

import "fmt"

// AlphanumericData header for storing string arrays
// case #3: alphanumeric data array header
type AlphanumericData struct {
	Flag         uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // BYTE      Always 2: Byte indicating an alphanumeric array.
	Filename     [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // WORD      Length of data following the header = length of string array + 3.
	UnusedByte   uint8    // BYTE      Unused byte.
	VariableName byte     // BYTE      = (1..26 meaning A$..Z$) + 192.
	UnusedWord   uint16   // WORD      = 32768.
	Checksum     uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

func (b AlphanumericData) BlockType() string {
	return "Header"
}

func (b AlphanumericData) Name() string {
	return "Alphanumeric Data Array"
}

// ToString returns a formatted string for the header
func (b AlphanumericData) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", b.BlockType(), b.Name())
	str += fmt.Sprintf("     - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("     - Variable Name: %c", b.VariableName-192)
	return str
}
