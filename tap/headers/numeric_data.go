package headers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/mrcook/tzxit/tape"
)

// NumericData header for storing numeric arrays.
// Case #2: numeric data array header.
type NumericData struct {
	Length uint16 // Length of the data in this block

	Flag         uint8    // BYTE     Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // BYTE     Always 1: Byte indicating a numeric array.
	Filename     [10]byte // CHAR[10] Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // WORD     Length of data following the header = length of number array * 5 + 3.
	UnusedByte   uint8    // BYTE     Unused byte.
	VariableName byte     // BYTE     = (1..26 meaning A..Z) + 128.
	UnusedWord   uint16   // WORD     = 32768.
	Checksum     uint8    // BYTE     Simply all bytes XORed (including flag byte).
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *NumericData) Read(reader *bufio.Reader) {
	// TODO: is fatal acceptable here?
	if length, err := tape.PeekBlockLength(reader); err != nil {
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

// ToString returns a formatted string for the header
func (b NumericData) ToString() string {
	str := fmt.Sprintf(" - Header       : %s\n", b.Name())
	str += fmt.Sprintf("     - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("     - Variable Name: %c", b.VariableName-128)
	return str
}
