package headers

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/mrcook/retroio/storage"
)

// ByteData header for storing Machine Code or Screens.
// Case #4: byte header or SCREEN$ header.
type ByteData struct {
	Length uint16 // Length of the data in this block

	Flag         uint8    // Always 0: byte indicating a standard ROM loading header.
	DataType     uint8    // Always 3: Byte indicating a byte header.
	ProgramName  [10]byte // Loading name of the program. Filled with spaces (0x20) to 10 characters.
	DataLength   uint16   // Length of data following the header, in case of a SCREEN$ header = 6912.
	StartAddress uint16   // In case of a SCREEN$ header = 16384.
	UnusedWord   uint16   // 32768.
	Checksum     uint8    // Simply all bytes XORed (including flag byte).
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (b *ByteData) Read(reader *storage.Reader) {
	// TODO: is fatal acceptable here?
	if length, err := reader.PeekShort(); err != nil {
		log.Fatalf("unexpected error reading block %v.", err)
	} else if length != 19 {
		log.Fatalf("expected header length to be 19, got '%d'.", length)
	}

	_ = binary.Read(reader, binary.LittleEndian, b)
}

func (b ByteData) Id() uint8 {
	return b.DataType
}

func (b ByteData) Name() string {
	if b.DataLength == 6912 && b.StartAddress == 16384 {
		return "SCREEN$"
	} else {
		return "Machine Code"
	}
}

func (b ByteData) Filename() string {
	return string(b.ProgramName[:])
}

func (b ByteData) BlockData() []byte {
	return []byte{}
}

// String returns a formatted string for the header
func (b ByteData) String() string {
	str := fmt.Sprintf("%s\n", b.Name())
	str += fmt.Sprintf("    - Filename     : %s\n", b.ProgramName)
	str += fmt.Sprintf("    - Start Address: %d", b.StartAddress)
	return str
}
