// Package tap implements reading of ZX Spectrum TAP formatted files,
// as specified in the TAP- and BLK specification at:
// http://www.zx-modules.de/fileformats/tapformat.html
//
// The TAP- (and BLK-) format is nearly a direct copy of the data that is
// stored in real tapes as is written by the ROM save routine of the
// ZX Spectrum.
//
// It is simply one data block, or a group of 2 or more data blocks, one
// followed after the other. A TAP file may be empty, but then it has a
// size of 0 bytes. There's no real file size limit, like real tapes, so
// they can contain a very large number of data blocks.
package tap

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"retroio/spectrum/tap/blocks"
	"retroio/spectrum/tap/headers"
	"retroio/storage"
)

// A TAP file may contain zero or more header/data block pairs.
type TAP struct {
	reader *storage.Reader

	Blocks []Block
}

// A Block as stored on tape may be a header or any data from the ZX Spectrum.
// NOTE: these blocks are also re-used in the TZX tape format.
//
// Blocks must have at least 2 bytes to build up a valid tape data block
// as each block should be bookended by `flag` and `checksum` bytes.
// Shorter data blocks (0 or 1 byte long) are so-called _fragmented_ ones.
type Block struct {
	Length   uint16
	TapeData BlockI
}

// Block is an interface for TAP header/data block
type BlockI interface {
	Read(reader *storage.Reader)
	Id() uint8
	Name() string
}

func New(reader *storage.Reader) *TAP {
	return &TAP{reader: reader}
}

// Read processes each TAP/BLK block in the tape file.
func (t *TAP) Read() error {
	// It's possible that a data block is 19 bytes long, but no two header blocks
	// can follow each other, so this check is required for those rare encounters
	// e.g. Turbo Outrun.
	blockCanBeHeader := true

	for {
		// Lookup the length of the block to know what type it is.
		blockLength, err := t.reader.PeekShort()
		if err != nil && err == io.EOF {
			break // no problems, we're done!
		} else if err != nil {
			return err
		}

		block := Block{Length: blockLength}

		if block.Length == 19 && blockCanBeHeader {
			block.TapeData, err = t.ReadHeaderBlock()
			blockCanBeHeader = false
		} else {
			block.TapeData, err = t.ReadDataBlock()
			blockCanBeHeader = true
		}

		if err != nil && err == io.EOF {
			break // no problems, we're done!
		} else if err != nil {
			return err
		}

		t.Blocks = append(t.Blocks, block)
	}

	return nil
}

// ReadHeaderBlock reads the different types of 19-byte header blocks.
func (t *TAP) ReadHeaderBlock() (BlockI, error) {
	// Look up the Flag and DataType bytes, ignoring the 2-byte block Length
	blockBytes, err := t.reader.Peek(4)
	if err != nil {
		return nil, err
	}
	// length := blockBytes[0:1]
	flag := blockBytes[2]
	dataType := blockBytes[3]

	if flag != 0 {
		return nil, errors.New(fmt.Sprintf("expected header FLAG byte to be 0, got '%d'", flag))
	}

	var header BlockI

	switch dataType {
	case 0:
		header = &headers.ProgramData{}
	case 1:
		header = &headers.NumericData{}
	case 2:
		header = &headers.AlphanumericData{}
	case 3:
		header = &headers.ByteData{}
	default:
		return nil, errors.New(fmt.Sprintf("unknown header type '%d'", dataType))
	}

	header.Read(t.reader)

	return header, nil
}

// ReadDataBlock reads a fragment or standard data block.
func (t *TAP) ReadDataBlock() (BlockI, error) {
	// Lookup the length of the block to know what type of block it is.
	length, err := t.reader.PeekShort()
	if err != nil {
		return nil, err
	}

	var block BlockI

	// Fragments are either 0 or 1 bytes long
	if length < 2 {
		block = &blocks.Fragment{}
	} else {
		block = &blocks.Standard{}
	}

	block.Read(t.reader)

	return block, nil
}

// DisplayImageMetadata outputs the metadata of each data block to the terminal.
func (t TAP) DisplayImageMetadata() {
	fmt.Println("TAP processing complete!")
	fmt.Println()
	fmt.Println("DATA BLOCKS:")
	for i, block := range t.Blocks {
		fmt.Printf("#%02d %s\n", i+1, block.TapeData)
	}
}