package headers

import (
	"encoding/binary"
	"fmt"
	"log"

	"retroio/storage"
)

// ProgramData header for storing BASIC programs.
// Case #1: program header or program auto-start header.
type ProgramData struct {
	Length uint16 // Length of the data in this block

	Flag          uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType      uint8    // BYTE      Always 0: Byte indicating a program header.
	Filename      [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength    uint16   // WORD      Length of data following the header = length of BASIC program + variables.
	AutoStartLine uint16   // WORD      LINE parameter of SAVE command. Value 32768 means "no auto-loading". 0..9999 are valid line numbers.
	ProgramLength uint16   // WORD      Length of BASIC program; remaining bytes ([data length] - [program length]) = offset of variables.
	Checksum      uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *ProgramData) Read(reader *storage.Reader) {
	// TODO: is fatal acceptable here?
	if length, err := reader.PeekShort(); err != nil {
		log.Fatalf("unexpected error reading block %v.", err)
	} else if length != 19 {
		log.Fatalf("expected header length to be 19, got '%d'.", length)
	}

	_ = binary.Read(reader, binary.LittleEndian, b)
}

func (b ProgramData) Id() uint8 {
	return b.DataType
}

func (b ProgramData) Name() string {
	return "BASIC Program"
}

// ToString returns a formatted string for the header
func (b ProgramData) ToString() string {
	str := fmt.Sprintf("%s\n", b.Name())
	str += fmt.Sprintf("    - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("    - AutoStartLine: %d", b.AutoStartLine)
	return str
}
