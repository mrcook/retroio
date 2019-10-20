package blocks

import (
	"bufio"
	"fmt"

	"mrcook/tzxit/tape"
)

// CswRecording
// ID: 18h (24d)
// This block contains a sequence of raw pulses encoded in CSW format v2 (Compressed Square Wave).
type CswRecording struct {
	Length           uint32  // DWORD    Block length (without these four bytes)
	Pause            uint16  // WORD     Pause after this block (in ms).
	SampleRate       uint16  // BYTE[3]  Sampling rate
	SampleSpareByte  uint8   // NOTE: `SampleRate` above uses only 2-bytes for value but specification says 3-bytes, so this is for the spare.
	CompressionType  uint8   // BYTE     Compression type: RLE, Z-RLE
	StoredPulseCount uint32  // DWORD    Number of stored pulses (after decompression, for validation purposes)
	Data             []uint8 // BYTE[N]  CSW data, encoded according to the CSW file format specification.
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CswRecording) Read(reader *bufio.Reader) {
	c.Length = tape.ReadLong(reader)
	c.Pause = tape.ReadShort(reader)
	c.SampleRate = tape.ReadShort(reader)
	c.SampleSpareByte, _ = reader.ReadByte()
	c.CompressionType, _ = reader.ReadByte()
	c.StoredPulseCount = tape.ReadLong(reader)

	// Yep, we're discarding the data for the moment
	tape.ReadNextBytes(reader, int(c.Length))
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (c CswRecording) Id() uint8 {
	return 0x18
}

// Name of the block as given in the TZX specification.
func (c CswRecording) Name() string {
	return "CSW Recording"
}

// ToString returns a human readable string of the block data
func (c CswRecording) ToString() string {
	str := fmt.Sprintf("> %s\n", c.Name())
	str += fmt.Sprintf(" - Pause (ms.): %d\n", c.Pause)
	str += fmt.Sprintf(" - Sample Rate: %d\n", c.SampleRate)
	str += fmt.Sprintf(" - Compression: %d\n", c.CompressionType)
	str += fmt.Sprintf(" - Pulse Count: %d\n", c.StoredPulseCount)

	return str
}
