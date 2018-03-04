package tape

import (
	"encoding/binary"
	"log"
	"os"
)

type File struct {
	file *os.File
}

func (f *File) Open(filename string) error {
	var err error
	f.file, err = os.Open(filename)
	if err != nil {
		return err
	}
	return nil
}

func (f File) Close() error {
	return f.file.Close()
}

// ReadByte a BYTE from the TZX.
// Note, this will often be used to get the block ID byte.
func (f File) ReadByte() (byte, error) {
	var b [1]byte
	_, err := f.file.Read(b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadShort reads a WORD from the TZX.
// The two bytes are processed as a _Little Endian_ value and returned as a uint16.
func (f *File) ReadShort() uint16 {
	b := f.ReadBytes(2)
	return binary.LittleEndian.Uint16(b[:])
}

// ReadSignedShort reads a `signed` WORD from the TZX.
// The two bytes are processed as a _Little Endian_ value and returned as a int16.
func (f *File) ReadSignedShort() int16 {
	b := f.ReadBytes(2)
	return int16(binary.LittleEndian.Uint16(b[:]))
}

// ReadLong reads a DWORD from the TZX.
// The four bytes are processed as a _Little Endian_ value and returned as a uint32.
func (f *File) ReadLong() uint32 {
	b := f.ReadBytes(4)
	return binary.LittleEndian.Uint32(b[:])
}

// ReadBytes reads a variable length of bytes from the TZX.
// Note: expected only to be used for reading block data so we should
// never reach EOF - if that happens then there's something seriously
// wrong with the tape image - so exit with a fatal.
func (f File) ReadBytes(number int) []byte {
	b := make([]byte, number)
	_, err := f.file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// BytesToLong converts a slice of bytes to a DWORD value
func (f *File) BytesToLong(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}
