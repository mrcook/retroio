package tap

import "fmt"

// ByteHeader for storing Machine Code or Screens
// case #4: byte header or SCREEN$ header
type ByteHeader struct {
	Flag         uint8    // BYTE      flag byte      always 0. Byte indicating a standard ROM loading header
	DataType     uint8    // BYTE      data type      always 3: Byte indicating a byte header
	Filename     [10]byte // CHAR[10]  file name      loading name of the program. filled with spaces (CHR$(32))
	DataLength   uint16   // WORD      [data length]  length of the following data (after the header) in case of a SCREEN$ header = 6912
	StartAddress uint16   // WORD      start address  in case of a SCREEN$ header = 16384
	UnusedWord   uint16   // WORD      unused         = 32768
	Checksum     uint8    // BYTE      checksum byte  simply all bytes (including flag byte) XORed
}

// ToString returns a formatted string for the header
func (h ByteHeader) ToString() string {
	str := fmt.Sprintf("  - %-13s: %s\n", flagLabel(h.Flag), typeLabel(h.DataType))
	str += fmt.Sprintf("  - Filename     : %s\n", h.Filename)
	str += fmt.Sprintf("  - Start Address: %d", h.StartAddress)
	return str
}
