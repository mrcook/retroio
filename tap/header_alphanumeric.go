package tap

import "fmt"

// AlphanumericHeader for storing string arrays
// case #3: alphanumeric data array header
type AlphanumericHeader struct {
	Flag         uint8    // BYTE      flag byte      always 0. Byte indicating a standard ROM loading header
	DataType     uint8    // BYTE      data type      always 2: Byte indicating an alphanumeric array
	Filename     [10]byte // CHAR[10]  file name      loading name of the program. filled with spaces (CHR$(32))
	DataLength   uint16   // WORD      [data length]  length of the following data (after the header) = length of string array +3
	UnusedByte   uint8    // BYTE      unused         unused byte
	VariableName byte     // BYTE      variable name  = (1..26 meaning A$..Z$) +192
	UnusedWord   uint16   // WORD      unused         = 32768
	Checksum     uint8    // BYTE      checksum byte  simply all bytes (including flag byte) XORed
}

// ToString returns a formatted string for the header
func (h AlphanumericHeader) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", flagLabel(h.Flag), typeLabel(h.DataType))
	str += fmt.Sprintf("     - Filename     : %s\n", h.Filename)
	str += fmt.Sprintf("     - Variable Name: %c", h.VariableName-192)
	return str
}
