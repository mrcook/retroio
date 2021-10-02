package blocks

import (
	"log"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// C64 ROM Type Data Block (deprecated in v1.20)
// ID: 16h (22d)
//
// Well, this block was created to support the Commodore 64 standard ROM and
// similar tape blocks. It is made so basically anything that uses two or four
// pulses (which are the same in pairs) per bit can be written with it. Some
// explanation:
//
// * A wave consists of TWO pulses. The structure contains the length of ONE pulse.
// * The wave MUST always start with the LOW amplitude, since the C64 can only
//   detect the transition HIGH -> LOW.
// * If some pulse length is 0 then the whole wave must not be present.
//   This applies to DATA too.
// * The XOR checksum (if it is set to 0 or 1) is a XOR of all bits inthe byte
//   XOR-ed with the value in this field as the start value.
// * Finish Byte waves should be played after each byte EXCEPT last one.
// * Finish Data waves should be ONLY played after last byte of data.
// * When all the Data has finished there is an optional Trailer Tone, which is
//   standard for the Repeated Blocks in C64 ROM Loader.
//
// The replay procedure looks like this:
//
// 1. Pilot Tone
// 2. Sync waves
// 3. Data Bytes (with XOR and/or Finish Byte waves)
// 4. Finish Data pulses
// 5. Trailing Tone
//
// NOTE: the numbers in brackets `[]` represent the values for C64 ROM loader.
type C64RomTypeData struct {
	blockID types.BlockType

	totalBlockLength        uint32 // Length of the WHOLE block including the data (extension rule)
	pilotToneLength         uint16 // PILOT TONE pulse length [616]
	numWavesInPilotTone     uint16 // Number of waves in PILOT TONE
	syncWavePulse1Length    uint16 // SYNC first wave pulse length [1176]
	syncWavePulse2Length    uint16 // SYNC 2nd wave pulse length [896]
	zeroBitWavePulse1Length uint16 // ZERO bit 1st wave pulse length [616]
	zeroBitWavePulse2Length uint16 // ZERO bit 2nd wave pulse length [896]
	oneBitWavePulse1Length  uint16 // ONE bit 1st wave pulse length [896]
	oneBitWavePulse2Length  uint16 // ONE bit 2nd wave pulse length [616]

	// XOR Checksum bit for each Data byte: [1]
	//   00 - Start XOR checksum with value 0
	//   01 - Start XOR checksum with value 1
	//   FF - No checksum bit
	xorChecksum uint8

	finishByteWavePulse1Length uint16 // FINISH BYTE 1st wave pulse length [1176]
	finishByteWavePulse2Length uint16 // FINISH BYTE 2nd wave pulse length [896]
	finishDataWavePulse1Length uint16 // FINISH DATA 1st wave pulse length [1176]
	finishDataWavePulse2Length uint16 // FINISH DATA 2nd wave pulse length [616]
	trailingTonePulseLength    uint16 // TRAILING TONE pulse length [616]
	numWavesInTrailingTone     uint16 // Number of waves in TRAILING TONE

	// Used bits in last byte (other bits should be 0)
	// e.g. if this is 6, then the bits used (x) in last byte are: xxxxxx00
	usedBitsInLastByte uint8

	// General Purpose, bit-mapped: [1]
	//   bit 0 - Data Endian-ness:
	//     0 = LSb first (little endian)
	//     1 = MSb first (big endian)
	endianFormat uint8

	pause  uint16        // Pause after this block in milliseconds (ms.)
	length [3]uint8      // Length (24-bit) of following data
	data   tap.TapeBlock // Data as in .TAP files
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (d *C64RomTypeData) Read(reader *storage.Reader) error {
	log.Fatal("C64RomTypeData is not currently supported - unable to continue.")

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (d C64RomTypeData) Id() types.BlockType {
	return types.C64RomType
}

// Name of the block as given in the TZX specification.
func (d C64RomTypeData) Name() string {
	return "C64 ROM Type Data"
}

func (d C64RomTypeData) BlockData() tap.Block {
	return nil
}
