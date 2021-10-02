package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/csw2"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// CswRecording
// ID: 18h (24d)
// This block contains a sequence of raw pulses encoded in CSW format v2 (Compressed Square Wave).
type CswRecording struct {
	BlockID          types.BlockType
	Length           uint32   // Block length (without these four bytes)
	Pause            uint16   // Pause after this block (in ms).
	SampleRate       uint16   // Sampling rate
	SampleSpareByte  uint8    // NOTE: `SampleRate` above uses only 2-bytes for value but specification says 3-bytes, so this is for the spare.
	CompressionType  uint8    // Compression type: RLE, Z-RLE
	StoredPulseCount uint32   // Number of stored pulses (after decompression, for validation purposes)
	Data             csw2.CSW // CSW data, encoded according to the CSW file format specification.
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CswRecording) Read(reader *storage.Reader) error {
	c.BlockID = types.BlockType(reader.ReadByte())
	if c.BlockID != c.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", c.Id(), c.BlockID)
	}

	c.Length = reader.ReadLong()
	c.Pause = reader.ReadShort()
	c.SampleRate = reader.ReadShort()
	c.SampleSpareByte = reader.ReadByte()
	c.CompressionType = reader.ReadByte()
	c.StoredPulseCount = reader.ReadLong()

	c.Data = csw2.CSW{}
	if err := c.Data.Read(reader, c.cswDataLength()); err != nil {
		return err
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (c CswRecording) Id() types.BlockType {
	return types.CswRecording
}

// Name of the block as given in the TZX specification.
func (c CswRecording) Name() string {
	return "CSW Recording"
}

func (c CswRecording) BlockData() tap.Block {
	return nil
}

// Length of this block size, without the length and CSW2 data
// This is required to know how much CSW data to read.
const CswRecordingHeaderSize uint32 = 10

func (c CswRecording) cswDataLength() int {
	return int(c.Length - CswRecordingHeaderSize)
}
