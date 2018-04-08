package tap

import "fmt"

// ProgramHeader for storing BASIC programs
// case #1: program header or program auto-start header
type ProgramHeader struct {
	Flag          uint8    // BYTE      flag byte         always 0. Byte indicating a standard ROM loading header
	DataType      uint8    // BYTE      data type         always 0: Byte indicating a program header
	Filename      [10]byte // CHAR[10]  file name         loading name of the program. filled with spaces (CHR$(32))
	DataLength    uint16   // WORD      [data length]     length of the following data (after the header) = length of BASIC program + variables
	AutoStartLine uint16   // WORD      autostart line    LINE parameter of SAVE command. Value 32768 means "no auto-loading"; 0..9999 are valid line numbers.
	ProgramLength uint16   // WORD      [program length]  length of BASIC program; remaining bytes ([data length] - [program length]) = offset of variables
	Checksum      uint8    // BYTE      checksum byte     simply all bytes (including flag byte) XORed
}

// ToString returns a formatted string for the header
func (h ProgramHeader) ToString() string {
	str := fmt.Sprintf("     - %-13s: %s\n", flagLabel(h.Flag), typeLabel(h.DataType))
	str += fmt.Sprintf("     - Filename     : %s\n", h.Filename)
	str += fmt.Sprintf("     - AutoStartLine: %d", h.AutoStartLine)
	return str
}
