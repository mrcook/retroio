package headers

import (
	"encoding/binary"
	"fmt"
	"log"

	"retroio/storage"
)

// NumericData header for storing numeric arrays.
// Case #2: numeric data array header.
type NumericData struct {
	Length uint16 // Length of the data in this block

	Flag         uint8    // Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // Always 1: Byte indicating a numeric array.
	ProgramName  [10]byte // Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // Length of data following the header = length of number array * 5 + 3.
	UnusedByte   uint8    // Unused byte.
	VariableName byte     // (1..26 meaning A..Z) + 128.
	UnusedWord   uint16   // 32768.
	Checksum     uint8    // Simply all bytes XORed (including flag byte).
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *NumericData) Read(reader *storage.Reader) {
	// TODO: is fatal acceptable here?
	if length, err := reader.PeekShort(); err != nil {
		log.Fatalf("unexpected error reading block %v.", err)
	} else if length != 19 {
		log.Fatalf("expected header length to be 19, got '%d'.", length)
	}

	_ = binary.Read(reader, binary.LittleEndian, b)
}

func (b NumericData) Id() uint8 {
	return b.DataType
}

func (b NumericData) Name() string {
	return "Numeric Data Array"
}

func (b NumericData) Filename() string {
	return string(b.ProgramName[:])
}

func (b NumericData) BlockData() []byte {
	return []byte{}
}

// String returns a formatted string for the header
func (b NumericData) String() string {
	str := fmt.Sprintf("%s\n", b.Name())
	str += fmt.Sprintf("    - Filename     : %s\n", b.ProgramName)
	str += fmt.Sprintf("    - Variable Name: %c", b.VariableName-128)
	return str
}
