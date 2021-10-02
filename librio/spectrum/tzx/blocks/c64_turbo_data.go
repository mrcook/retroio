package blocks

import (
	"log"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// C64 Turbo Tape Data Block (deprecated in v1.20)
// ID: 17h (23d)
//
// This block is made to support another type of encoding that is commonly
// used by the C64. Most of the commercial software uses this type of
// encoding, i.e. the Pilot tone is not made from one type of Wave only,
// but it is made from actual Data byte which is repeated many times. As
// the Sync value another, different, Data byte is sent to signal the
// start of the data. The Data Bits are made from ONE wave only and there
// is NO XOR checksum either! Trailing byte is played AFTER the DATA has
// ended.
//
// NOTE: the numbers in brackets `[]` represent the values for C64 ROM loader.
type C64TurboData struct {
	blockID types.BlockType

	totalBlockLength uint32 // Length of the WHOLE block including the data (extension rule)
	zeroBitPulse     uint16 // ZERO bit pulse
	oneBitPulse      uint16 // ONE bit pulse

	// Additional bits in bytes (bit-mapped)
	//   bits 0-1: number of bits (0-3)
	//   bit    2: play additional bit(s)
	//             BEFORE (0) or AFTER (1) the byte
	//   bit    3: value of additional bit(s)
	additionalBits uint8

	numLeadInBytes uint16 // Number of lead-in bytes
	leadInByte     uint8  // Lead-in byte

	// Used bits in last byte (other bits should be 0) [8]
	// e.g. if this is 6, then the bits used (x) in last byte are: xxxxxx00
	usedBitsInLastByte uint8

	// General Purpose, bit-mapped: [0]
	//   bit 0 - Data Endian-ness:
	//     0 = LSb first (little endian)
	//     1 = MSb first (big endian)
	endianFormat uint8

	numTrailingBytes uint16        // Number of trailing bytes
	trailingByte     uint8         // Trailing byte
	pause            uint16        // Pause after this block in milliseconds (ms.)
	length           [3]uint8      // Length (24-bit) of following data
	data             tap.TapeBlock // Data as in .TAP files
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (d *C64TurboData) Read(reader *storage.Reader) error {
	log.Fatal("C64TurboData is not currently supported - unable to continue.")

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (d C64TurboData) Id() types.BlockType {
	return types.C64TurboData
}

// Name of the block as given in the TZX specification.
func (d C64TurboData) Name() string {
	return "C64 Turbo Data"
}

func (d C64TurboData) BlockData() tap.Block {
	return nil
}
