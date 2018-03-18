// Package tap processes the ZX Spectrum TAP/BLK header/data blocks
package tap

// The TAP- (and BLK-) format is nearly a direct copy of the data that is
// stored in real tapes, as it is 'written' by the ROM save routine of the
// ZX-Spectrum. It is supported by many Spectrum emulators.

// The TAP "tape data blocks" are also used for the TZX file format.

// A TAP file is simply one "data block" or a group of 2 or more data blocks,
// one followed after the other. The TAP file may be empty, but then it has a
// size of 0 bytes. There's no real file size limit, like real tapes, TAP files
// can also contain huge amounts of data blocks.

// Data blocks must have at least 2 bytes, to build up a valid tape data block,
// because each block should have a `flag` and a `checksum` byte.

// Shorter data blocks (1 or 0 bytes long) are so-called "fragmented" ones.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Block provides an interface to TAP blocks
type Block interface {
	ToString() string
}

// header to ease unmarshalling
type header struct {
	flag byte
	typ  byte
	data []byte
}

// Unmarshal a 19 byte data block onto a TAP header
func Unmarshal(data []byte) (Block, error) {
	if len(data) != 19 {
		return nil, fmt.Errorf("incorrect data length")
	}

	h := header{data[0], data[1], data}

	var block Block
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
func (h *header) unmarshalRomHeader() (Block, error) {
	switch h.typ {
	case 0:
		block := ProgramHeader{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 1:
		block := NumericHeader{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 2:
		block := AlphanumericHeader{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	case 3:
		block := ByteHeader{}
		buffer := bytes.NewBuffer(h.data)
		err := binary.Read(buffer, binary.LittleEndian, &block)
		return block, err
	default:
		return nil, fmt.Errorf("invalid data type byte: %d", h.typ)
	}
}

// unmarshalRomBlock handles ROM loading blocks
func (h *header) unmarshalRomBlock() (Block, error) {
	d := StandardDataBlocks{
		Flag:     h.flag,
		Data:     h.data[1:18],
		Checksum: h.data[19],
	}
	return d, nil
}

// unmarshalCustomBlock handles custom data blocks
func (h *header) unmarshalCustomBlock() (Block, error) {
	d := StandardDataBlocks{
		Flag:     h.flag,
		Data:     h.data[1:18],
		Checksum: h.data[19],
	}
	return d, nil
}

// flagLabel returns a human readable header `flag` label
func flagLabel(flag byte) string {
	switch flag {
	case 0:
		return "Header"
	case 255:
		return "Block"
	default:
		return "Custom Bock"
	}
}

// typeLabel returns a human readable header `type` label
func typeLabel(typ byte) string {
	switch typ {
	case 0:
		return "BASIC Program"
	case 1:
		return "Numeric Array"
	case 2:
		return "Alphanumeric Array"
	case 3:
		return "Machine Code/Screen"
	default:
		return "WARNING: unknown type"
	}
}
