package headers

import "fmt"

// ProgramData header for storing BASIC programs
// case #1: program header or program auto-start header
type ProgramData struct {
	Flag          uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType      uint8    // BYTE      Always 0: Byte indicating a program header.
	Filename      [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength    uint16   // WORD      Length of data following the header = length of BASIC program + variables.
	AutoStartLine uint16   // WORD      LINE parameter of SAVE command. Value 32768 means "no auto-loading". 0..9999 are valid line numbers.
	ProgramLength uint16   // WORD      Length of BASIC program; remaining bytes ([data length] - [program length]) = offset of variables.
	Checksum      uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

func (b ProgramData) BlockType() string {
	return "Header"
}

func (b ProgramData) Name() string {
	return "BASIC Program"
}

// ToString returns a formatted string for the header
func (b ProgramData) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", b.BlockType(), b.Name())
	str += fmt.Sprintf("     - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("     - AutoStartLine: %d", b.AutoStartLine)
	return str
}
