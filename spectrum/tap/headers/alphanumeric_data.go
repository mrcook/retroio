package headers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"

	"retroio/tape"
)

// AlphanumericData header for storing string arrays.
// Case #3: alphanumeric data array header.
type AlphanumericData struct {
	Length uint16 // Length of the data in this block

	Flag         uint8    // BYTE      Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // BYTE      Always 2: Byte indicating an alphanumeric array.
	Filename     [10]byte // CHAR[10]  Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // WORD      Length of data following the header = length of string array + 3.
	UnusedByte   uint8    // BYTE      Unused byte.
	VariableName byte     // BYTE      = (1..26 meaning A$..Z$) + 192.
	UnusedWord   uint16   // WORD      = 32768.
	Checksum     uint8    // BYTE      Simply all bytes XORed (including flag byte).
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *AlphanumericData) Read(reader *bufio.Reader) {
	// TODO: is fatal acceptable here?
	if length, err := tape.PeekBlockLength(reader); err != nil {
		log.Fatalf("unexpected error reading block %v.", err)
	} else if length != 19 {
		log.Fatalf("expected header length to be 19, got '%d'.", length)
	}

	_ = binary.Read(reader, binary.LittleEndian, b)
}

func (b AlphanumericData) Id() uint8 {
	return b.DataType
}

func (b AlphanumericData) Name() string {
	return "Alphanumeric Data Array"
}

// ToString returns a formatted string for the header
func (b AlphanumericData) ToString() string {
	str := fmt.Sprintf(" - Header       : %s\n", b.Name())
	str += fmt.Sprintf("     - Filename     : %s\n", b.Filename)
	str += fmt.Sprintf("     - Variable Name: %c", b.VariableName-192)
	return str
}
