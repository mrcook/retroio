// Package tape provides interfaces and byte reading helper functions for
// processing ZX Spectrum tape files.
//
// Tape files may contain different cassette tape file formats used by various
// emulators. Two of the most popular are TZX and TAP.
//   TZX tapes contain a metadata header, followed by one of more TZX blocks.
//   TAP tapes contain one of more data blocks.
package tape

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type Tape interface {
	ReadBlocks() error
	DisplayTapeMetadata()
}

// Block is an interface for Tape data blocks
type Block interface {
	Read(reader *bufio.Reader)
	Id() uint8
	Name() string
	ToString() string
}

// PeekBlockLength reads two bytes without advancing the reader,
// then converts the little endian ordered bytes to an integer.
func PeekBlockLength(r *bufio.Reader) (uint16, error) {
	b, err := r.Peek(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

// ReadNextBytes reads a variable length of bytes from the reader.
// NOTE: this is expected only to be used for reading block data, so it should
// never reach EOF. If that happens then there's something seriously wrong with
// the tape image, so a fatal is triggered.
func ReadNextBytes(r *bufio.Reader, number int) []byte {
	b := make([]byte, number)
	_, err := io.ReadFull(r, b)
	if err != nil {
		fmt.Println("tape.ReadNextBytes failed spectacularly!")
		log.Fatal(err)
	}
	return b
}

// ReadLong reads a DWORD from the reader, then converts the little endian
// ordered bytes to a uint32.
func ReadLong(r *bufio.Reader) uint32 {
	b := ReadNextBytes(r, 4)
	return binary.LittleEndian.Uint32(b[:])
}

// ReadShort reads a WORD from the reader, then converts the little endian
// ordered bytes to a uint16.
func ReadShort(r *bufio.Reader) uint16 {
	b := ReadNextBytes(r, 2)
	return binary.LittleEndian.Uint16(b[:])
}

// ReadSignedShort reads a `Signed` WORD from the reader, then converts the little endian
// ordered bytes to a int16.
func ReadSignedShort(r *bufio.Reader) int16 {
	b := ReadNextBytes(r, 2)
	return int16(binary.LittleEndian.Uint16(b[:]))
}

// BytesToLong converts a slice of bytes to a DWORD value.
func BytesToLong(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}
