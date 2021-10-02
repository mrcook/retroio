// The TZX library implements the cassette tape file format as described in the
// ZX Spectrum TZX specification. https://www.worldofspectrum.org/TZXformat.html
//
// Rules and Definitions
//
//  * Any value requiring more than one byte is stored in little endian format (i.e. LSB first).
//  * Unused bits should be set to zero.
//  * Timings are given in Z80 clock ticks (T states) unless otherwise stated.
//      1 T state = (1/3500000)s
//  * Block IDs are given in hex.
//  * All ASCII texts use the ISO 8859-1 (Latin 1) encoding; some of them can have several lines, which
//    should be separated by ASCII code 13 decimal (0D hex).
//  * You might interpret 'full-period' as ----____ or ____----, and 'half-period' as ---- or ____.
//    One 'half-period' will also be referred to as a 'pulse'.
//  * Values in curly brackets {} are the default values that are used in the Spectrum ROM saving
//    routines. These values are in decimal.
//  * If there is no pause between two data blocks then the second one should follow immediately; not
//    even so much as one T state between them.
//  * This document refers to 'high' and 'low' pulse levels. Whether this is implemented as ear=1 and
//    ear=0 respectively or the other way around is not important, as long as it is done consistently.
//  * Zeros and ones in 'Direct recording' blocks mean low and high pulse levels respectively.
//    The 'current pulse level' after playing a Direct Recording block of CSW recording block
//    is the last level played.
//  * The 'current pulse level' after playing the blocks ID 10,11,12,13,14 or 19 is the opposite of
//    the last pulse level played, so that a subsequent pulse will produce an edge.
//  * A 'Pause' block consists of a 'low' pulse level of some duration. To ensure that the last edge
//    produced is properly finished there should be at least 1 ms. pause of the opposite level and only
//    after that the pulse should go to 'low'. At the end of a 'Pause' block the 'current pulse level'
//    is low (note that the first pulse will therefore not immediately produce an edge). A 'Pause' block
//    of zero duration is completely ignored, so the 'current pulse level' will NOT change in this case.
//    This also applies to 'Data' blocks that have some pause duration included in them.
//  * An emulator should put the 'current pulse level' to 'low' when starting to play a TZX file, either
//    from the start or from a certain position. The writer of a TZX file should ensure that the 'current
//    pulse level' is well-defined in every sequence of blocks where this is important, i.e. in any
//    sequence that includes a 'Direct recording' block, or that depends on edges generated by 'Pause'
//    blocks. The recommended way of doing this is to include a Pause after each sequence of blocks.
//  * When creating a 'Direct recording' block please stick to the standard sampling frequencies of 22050
//    or 44100 Hz. This will ensure correct playback when using PC's sound cards.
//  * The length of a block is given in the following format: numbers in square brackets [] mean that the
//    value must be read from the offset in the brackets. Other values are normal numbers.
//    Example: [02,03]+0A means: get number (a word) from offset 02 and add 0A. All numbers are in hex.
//  * General Extension Rule: ALL custom blocks that will be added after version 1.10 will have the length
//    of the block in first 4 bytes (long word) after the ID (this length does not include these 4 length
//    bytes). This should enable programs that can only handle older versions to skip that block.
//  * Just in case:
//      MSB = most significant byte
//      LSB = least significant byte
//      MSb = most significant bit
//      LSb = least significant bit
package tzx

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks"

	"github.com/pkg/errors"

	"github.com/mrcook/retroio/librio/storage"
)

// TZX files store the header information at the start of the file, followed
// by zero or more data blocks. Some TZX files include an ArchiveInfo block,
// which is always stored as the first block, directly after the header.
type TZX struct {
	reader *storage.Reader

	Header Header

	// NOTE: Archive Info is always the first entry when present.
	Blocks []blocks.Block // TODO: needs to be a generic type
}

// Read processes the header, and then each block on the tape.
func (t *TZX) Read() error {
	if err := t.readTapeHeader(); err != nil {
		return err
	}

	if err := t.readBlocks(); err != nil {
		return err
	}

	return nil
}

// readHeader reads the tape header data and validates that the format is correct.
func (t *TZX) readTapeHeader() error {
	t.Header = Header{}

	if err := binary.Read(t.reader, binary.LittleEndian, &t.Header); err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	if err := t.Header.valid(); err != nil {
		return err
	}

	return nil
}

// readBlocks processes each TZX block on the tape.
func (t *TZX) readBlocks() error {
	for {
		blockID, err := t.reader.PeekByte()
		if err != nil {
			if err == io.EOF {
				break // no problems, we're done!
			}
			return err
		}

		block, err := blocks.NewFromBlockID(blockID)
		if err != nil {
			return err
		}

		if err := block.Read(t.reader); err != nil {
			return errors.Wrap(err, "error reading TZX block")
		}
		t.Blocks = append(t.Blocks, block)
	}

	return nil
}
