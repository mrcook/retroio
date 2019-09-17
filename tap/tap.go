// Package tap processes the ZX Spectrum TAP/BLK header/data blocks
package tap

// A TAP file is simply one _data block_ or a group of 2 or more data blocks,
// one followed after the other. The TAP file may be empty, but then it has a
// size of 0 bytes. There's no real file size limit, like real tapes, TAP files
// can also contain huge amounts of data blocks.

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mrcook/tzxit/tap/blocks"
	"github.com/mrcook/tzxit/tap/headers"
	"github.com/mrcook/tzxit/tape"
)

// The TAP- (and BLK-) format is nearly a direct copy of the data that is
// stored in real tapes as it is _written_ by the ROM save routine of the
// ZX-Spectrum.
// It's simply one _data block_ or a group of 2 or more data blocks, one
// followed after the other.
type Tap struct {
	Blocks []DataBlock
}

// A DataBlock as stored on tape: may be a header or any data from ZX-Spectrum.
// (Please note that these are re-used for the TZX file format as well)
//
// Data block must have at least 2 bytes to build up a valid tape data block,
// because each block should have a `flag` and a `checksum` byte.
// Shorter data blocks (1 or 0 bytes long) are so-called "fragmented" ones.
type DataBlock struct {
	Length uint16
	Data   interface{}
}

// Block provides an interface to TAP blocks
type BlockI interface {
	ToString() string
}

// Read the tape data blocks
func (t *Tap) Read(file *tape.Reader) error {
	return nil
}

// header to ease unmarshalling
type header struct {
	flag byte
	typ  byte
	data []byte
}

// Unmarshal a 19 byte data block onto a TAP header
func Unmarshal(data []byte) (BlockI, error) {
	if len(data) != 19 {
		return nil, fmt.Errorf("incorrect data length")
	}

	h := header{data[0], data[1], data}

	var block BlockI
	var err error

	switch h.flag {
	case 0:
		block, err = h.unmarshalRomHeader()
	case 255:
		block, err = h.unmarshalRomBlock()
	default:
		block, err = h.unmarshalCustomBlock()
	}

	if err != nil {
		return nil, err
	}

	return block, nil
}

// unmarshalRomHeader handles ROM loading headers
func (h *header) unmarshalRomHeader() (BlockI, error) {
	switch h.typ {
	case 0:
		block := headers.ProgramData{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 1:
		block := headers.NumericData{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 2:
		block := headers.AlphanumericData{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 3:
		block := headers.ByteData{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	default:
		return nil, fmt.Errorf("invalid data type byte: %d", h.typ)
	}
}

// unmarshalRomBlock handles ROM loading blocks
func (h *header) unmarshalRomBlock() (BlockI, error) {
	d := blocks.Standard{
		Flag:     h.flag,
		Data:     h.data[1:18],
		Checksum: h.data[19],
	}
	return d, nil
}

// unmarshalCustomBlock handles custom data blocks
func (h *header) unmarshalCustomBlock() (BlockI, error) {
	d := blocks.Standard{
		Flag:     h.flag,
		Data:     h.data[1:18],
		Checksum: h.data[19],
	}
	return d, nil
}
