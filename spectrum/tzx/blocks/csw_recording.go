package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// CswRecording
// ID: 18h (24d)
// This block contains a sequence of raw pulses encoded in CSW format v2 (Compressed Square Wave).
type CswRecording struct {
	BlockID          types.BlockType
	Length           uint32  // Block length (without these four bytes)
	Pause            uint16  // Pause after this block (in ms).
	SampleRate       uint16  // Sampling rate
	SampleSpareByte  uint8   // NOTE: `SampleRate` above uses only 2-bytes for value but specification says 3-bytes, so this is for the spare.
	CompressionType  uint8   // Compression type: RLE, Z-RLE
	StoredPulseCount uint32  // Number of stored pulses (after decompression, for validation purposes)
	Data             []uint8 // CSW data, encoded according to the CSW file format specification.
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

	c.Data = make([]byte, c.Length)
	if _, err := reader.Read(c.Data); err != nil {
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

// String returns a human readable string of the block data
func (c CswRecording) String() string {
	str := fmt.Sprintf("%s\n", c.Name())
	str += fmt.Sprintf(" - Pause (ms.): %d\n", c.Pause)
	str += fmt.Sprintf(" - Sample Rate: %d\n", c.SampleRate)
	str += fmt.Sprintf(" - Compression: %d\n", c.CompressionType)
	str += fmt.Sprintf(" - Pulse Count: %d\n", c.StoredPulseCount)

	return str
}
