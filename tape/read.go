package tape

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// A Reader is an io.Reader that can be read to retrieve the binary
// data from a file. It also includes useful function for reading
// specific int types.
//
// Tape files may contain different cassette tape file formats used by various
// emulators. Two of the most popular are TZX and TAP.
//   TZX files contain a header, followed by one of more data blocks.
//   TAP files contain one of more data blocks.
type Reader struct {
	io.Reader
}

// NewReader creates a new Reader reading the given reader.
//
// It is the caller's responsibility to call Close on the Reader when done.
func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

// ReadByte a BYTE from the tape.
// Note, this will often be used to get the block ID byte.
func (r Reader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadShort reads a WORD from the tape.
// The two bytes are processed as a _Little Endian_ value and returned as a uint16.
func (r Reader) ReadShort() uint16 {
	b := r.ReadBytes(2)
	return binary.LittleEndian.Uint16(b[:])
}

// ReadSignedShort reads a `signed` WORD from the tape.
// The two bytes are processed as a _Little Endian_ value and returned as a int16.
func (r Reader) ReadSignedShort() int16 {
	b := r.ReadBytes(2)
	return int16(binary.LittleEndian.Uint16(b[:]))
}

// ReadLong reads a DWORD from the tape.
// The four bytes are processed as a _Little Endian_ value and returned as a uint32.
func (r Reader) ReadLong() uint32 {
	b := r.ReadBytes(4)
	return binary.LittleEndian.Uint32(b[:])
}

// ReadBytes reads a variable length of bytes from the tape.
// Note: expected only to be used for reading block data so we should
// never reach EOF - if that happens then there's something seriously
// wrong with the tape image - so exit with a fatal.
func (r Reader) ReadBytes(number int) []byte {
	b := make([]byte, number)
	_, err := r.Read(b)
	if err != nil {
		fmt.Println("tape.ReadBytes failed spectacularly!")
		log.Fatal(err)
	}
	return b
}

// BytesToLong converts a slice of bytes to a DWORD value
func (r Reader) BytesToLong(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}
