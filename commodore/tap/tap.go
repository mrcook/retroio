// Package tap implements reading of Commodore C64 tape files as specified at:
// https://ist.uwaterloo.ca/~schepers/formats/TAP.TXT
//
// Designed by Per Hakan Sundell (author of the CCS64 C64 emulator) in 1997,
// this format attempts to duplicate the data stored on a C64 cassette tape,
// bit for bit. Since it is simply a representation of the raw serial data
// from a tape, it should handle *any* custom tape loaders that exist.
//
// The TAP images are generally very large, being a minimum of eight times,
// and up to sixteen times as large as what a raw PRG file would be. This is
// due to the way the data is stored, with each bit of the original file now
// being one byte large in the TAP file. The layout is fairly simple, with a
// small 14-byte header followed by file data.
//
// Note: DWORD values in the header are stored in little endian byte order.
package tap

import (
	"fmt"
	"io/ioutil"

	"retroio/storage"
)

// TAP File structure
type TAP struct {
	reader *storage.Reader

	Signature [12]byte // File signature "C64-TAPE-RAW"
	Version   uint8    // TAP version: $00 original layout, $01 updated.
	Unused    [3]byte  // Future expansion
	DataSize  uint32   // File data size (not including this header)
	Data      []byte   // File data: 0014-xxxx
}

func New(reader *storage.Reader) *TAP {
	return &TAP{reader: reader}
}

func (t *TAP) Read() error {
	if err := t.readTape(); err != nil {
		return err
	}

	return nil
}

// DisplayImageMetadata prints the tape metadata to the terminal.
func (t TAP) DisplayImageMetadata() {
	fmt.Println("HEADER INFORMATION:")
	fmt.Println(t.ToString())
}

func (t TAP) ToString() string {
	str := ""
	str += fmt.Sprintf("Signature  %s\n", t.Signature)
	str += fmt.Sprintf("Version:   $%02x (%s)\n", t.Version, t.tapType(t.Version))
	str += fmt.Sprintf("Data Size: %d bytes\n", t.DataSize)

	dataLenDiff := int(t.DataSize) - len(t.Data)
	if dataLenDiff != 0 {
		str += "\n"
		str += fmt.Sprintf("WARNING: data size mismatch, found %d bytes, %d difference\n", len(t.Data), dataLenDiff)
	}
	return str
}

// readTape reads the cassette tape header and data.
func (t *TAP) readTape() error {
	if _, err := t.reader.Read(t.Signature[:]); err != nil {
		return err
	}
	t.Version = t.reader.ReadByte()
	if _, err := t.reader.Read(t.Unused[:]); err != nil {
		return err
	}

	t.DataSize = t.reader.ReadLong()

	// Slurp the data
	var err error
	if t.Data, err = ioutil.ReadAll(t.reader); err != nil {
		return err
	}

	return nil
}

func (t TAP) tapType(id byte) string {
	var label string
	switch id {
	case 0x00:
		label = "Original Layout"
	case 0x01:
		label = "Updated Layout"
	default:
		label = "Unknown Layout"
	}
	return label
}
