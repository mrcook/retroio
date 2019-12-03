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
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"retroio/tap/blocks"
	"retroio/tap/headers"
	"retroio/tape"
)

// Reader wraps a bufio.Reader that can be used to read binary data from a tape
// file, but also provides addition functions for read TAP data blocks.
type Reader struct {
	reader *bufio.Reader
	Tap    *Tap
}

// A TAP file may contain zero or more data blocks.
type Tap struct {
	Blocks []DataBlock
}

// A DataBlock as stored on tape may be a header or any data from the ZX Spectrum.
// NOTE: these blocks are also re-used in the TZX tape format.
//
// Blocks must have at least 2 bytes to build up a valid tape data block
// as each block should be bookended by `flag` and `checksum` bytes.
// Shorter data blocks (0 or 1 byte long) are so-called _fragmented_ ones.
type DataBlock struct {
	Length   uint16
	TapeData tape.Block
}

// NewReader wraps the given buffered Reader and creates a new TAP Reader.
// NOTE: It's the caller's responsibility to call Close on the Reader when done.
func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		reader: r,
		Tap:    &Tap{},
	}
}

// ReadBlocks processes each TAP/BLK block in the tape file.
//
// NOTE: the TAP specification is a little unclear if each data block should be
// preceded by a header block. However, we will presume they do, therefore
// there are two consecutive reads; the header followed by the data block.
func (r *Reader) ReadBlocks() error {
	for {
		header, err := r.ReadHeaderBlock()
		if err != nil && err == io.EOF {
			break // no problems, we're done!
		} else if err != nil {
			return err
		}
		r.Tap.Blocks = append(r.Tap.Blocks, header)

		block, err := r.ReadDataBlock()
		if err != nil {
			// should never be an EOF error for valid tape files
			return errors.Wrap(err, "Unable to complete reading TZX blocks")
		}
		r.Tap.Blocks = append(r.Tap.Blocks, block)
	}
	return nil
}

// DisplayTapeMetadata outputs the metadata of each data block to the terminal.
func (r Reader) DisplayTapeMetadata() {
	fmt.Println("Tzxit TAP processing complete!")
	fmt.Println()
	fmt.Println("DATA BLOCKS:")
	for i, block := range r.Tap.Blocks {
		fmt.Printf("#%02d %s\n", i+1, block.TapeData.ToString())
	}
}

// ReadHeaderBlock reads the different types of header blocks.
func (r *Reader) ReadHeaderBlock() (DataBlock, error) {
	block := DataBlock{}

	length, err := tape.PeekBlockLength(r.reader)
	if err != nil {
		return block, err
	}

	// Look up the Flag and DataType bytes, ignoring the 2-byte block Length
	blockBytes, err := r.reader.Peek(4)
	if err != nil {
		return block, err
	}
	flag := blockBytes[2]
	dataType := blockBytes[3]

	if flag != 0 {
		return block, errors.New(fmt.Sprintf("expected header FLAG byte to be 0, got '%d'.", flag))
	}

	switch dataType {
	case 0:
		block.TapeData = &headers.ProgramData{}
	case 1:
		block.TapeData = &headers.NumericData{}
	case 2:
		block.TapeData = &headers.AlphanumericData{}
	case 3:
		block.TapeData = &headers.ByteData{}
	default:
		return block, errors.New(fmt.Sprintf("Unknown header type '%d'.", dataType))
	}

	block.Length = length
	block.TapeData.Read(r.reader)

	return block, nil
}

// ReadDataBlock reads a fragment or standard data block.
func (r *Reader) ReadDataBlock() (DataBlock, error) {
	var block DataBlock

	// Get the length of the block
	length, err := tape.PeekBlockLength(r.reader)
	if err != nil {
		return block, err
	}

	if length < 2 {
		block.TapeData = &blocks.Fragment{}
	} else {
		block.TapeData = &blocks.Standard{}
	}

	block.Length = length
	block.TapeData.Read(r.reader)

	return block, nil
}
