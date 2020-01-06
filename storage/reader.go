// Package image provides an io.Reader interface, along with various byte reading
// helper functions for processing image data files.
//
// To make this reader easier to use, most of these functions omit the handling
// of errors during the read operation. As we are working with known header/data
// structures, this has been done solely for ease of use.
// Should an EOF ever be reached...well...then there is something seriously wrong
// with the image file, so a `fatal` is triggered.
package storage

import (
	"bufio"
	"encoding/binary"
	"io"
)

// Image reader, using the bufio.Reader to allow for Peeking.
type Reader struct {
	reader *bufio.Reader
}

// NewReader first converts the regular reader to a buffered reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}

// Read exactly the requested bytes from the reader, and follows the reader interface.
// It will read either the currently buffered bytes, or perform a io.ReadFull.
func (r Reader) Read(b []byte) (int, error) {
	// if the buffer contains enough bytes, use them.
	if len(b) <= r.reader.Buffered() {
		return r.reader.Read(b)
	}

	return io.ReadFull(r.reader, b)
}

// ReadByte delegates to the underlying Reader function, and reads a single byte.
// Errors are discarded so this should only be used when a byte is known to be present.
func (r Reader) ReadByte() byte {
	b, _ := r.reader.ReadByte()
	return b
}

// ReadBytes reads a variable length of bytes from the reader.
// Errors are discarded so this should only be used when a byte is known to be present.
func (r Reader) ReadBytes(number int) []byte {
	b := make([]byte, number)
	_, _ = r.Read(b)
	return b
}

// ReadShort reads a value from the reader, converting the little endian ordered bytes to a uint16.
func (r Reader) ReadShort() uint16 {
	b := r.ReadBytes(2)
	return binary.LittleEndian.Uint16(b[:])
}

// ReadLong reads a value from the reader, converting the little endian ordered bytes to a uint32.
func (r Reader) ReadLong() uint32 {
	b := r.ReadBytes(4)
	return binary.LittleEndian.Uint32(b[:])
}

// Buffered delegates to the underlying Reader function, returning the number of bytes left in the buffer.
func (r Reader) Buffered() int {
	return r.reader.Buffered()
}

// Peek returns the next n bytes without advancing the reader.
func (r Reader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}

// PeekByte reads a byte without advancing the reader.
func (r Reader) PeekByte() (uint8, error) {
	b, err := r.reader.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// PeekShort reads two bytes without advancing the reader, converting
// the little endian ordered bytes to a uint16.
func (r Reader) PeekShort() (uint16, error) {
	b, err := r.reader.Peek(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

// Discard delegates to the underlying Reader function.
func (r Reader) Discard(n int) (int, error) {
	return r.reader.Discard(n)
}

// BytesToLong converts a slice of 4 little endian ordered bytes to uint32.
func (r Reader) BytesToLong(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}

// Bytes3ToLong converts a slice of 3 little endian ordered bytes to uint32.
func (r Reader) Bytes3ToLong(b [3]byte) uint32 {
	l := append(b[:], 0) // add 4th byte
	return binary.LittleEndian.Uint32(l[:])
}
