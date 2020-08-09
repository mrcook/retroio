package headers

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/mrcook/retroio/storage"
)

// ProgramData header for storing BASIC programs.
// Case #1: program header or program auto-start header.
type ProgramData struct {
	Length uint16 // Length of the data in this block

	Flag          uint8    // Always 0: byte indicating a standard ROM loading header.
	DataType      uint8    // Always 0: Byte indicating a program header.
	ProgramName   [10]byte // Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength    uint16   // Length of data following the header = length of BASIC program + variables.
	AutoStartLine uint16   // LINE parameter of SAVE command. Value 32768 means "no auto-loading". 0..9999 are valid line numbers.
	ProgramLength uint16   // Length of BASIC program; remaining bytes ([data length] - [program length]) = offset of variables.
	Checksum      uint8    // Simply all bytes XORed (including flag byte).
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

func (b ProgramData) Filename() string {
	return string(b.ProgramName[:])
}

func (b ProgramData) BlockData() []byte {
	return []byte{}
}

// String returns a formatted string for the header
func (b ProgramData) String() string {
	str := fmt.Sprintf("%s\n", b.Name())
	str += fmt.Sprintf("    - Filename        : %s\n", b.ProgramName)
	str += fmt.Sprintf("    - AutoStartLine   : %d", b.AutoStartLine)
	return str
}
