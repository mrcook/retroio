package headers

import "fmt"

// ByteData header for storing Machine Code or Screens
// case #4: byte header or SCREEN$ header
type ByteData struct {
	Flag         uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // BYTE      Always 3: Byte indicating a byte header.
	Filename     [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // WORD      Length of data following the header, in case of a SCREEN$ header = 6912.
	StartAddress uint16   // WORD      In case of a SCREEN$ header = 16384.
	UnusedWord   uint16   // WORD      = 32768.
	Checksum     uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

func (b ByteData) BlockType() string {
	return "Header"
}

func (b ByteData) Name() string {
	return "Byte Data: Machine Code/Screen"
}

// ToString returns a formatted string for the header
func (b ByteData) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", b.BlockType(), b.Name())
	str += fmt.Sprintf("     - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("     - Start Address: %d", b.StartAddress)
	return str
}
