package tape

import "fmt"

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

func (c *CswRecording) Process(file *File) {
	c.Length = file.ReadLong()
	c.Pause = file.ReadShort()
	c.SampleRate = file.ReadShort()
	c.SampleSpareByte, _ = file.ReadByte()
	c.CompressionType, _ = file.ReadByte()
	c.StoredPulseCount = file.ReadLong()

	// Yep, we're discarding the data for the moment
	file.ReadBytes(int(c.Length))
}

func (c CswRecording) Id() int {
	return 24
}

func (c CswRecording) Name() string {
	return "CSW Recording"
}

// Metadata returns a human readable string of the block data
func (c CswRecording) Metadata() string {
	str := fmt.Sprintf("> %s\n", c.Name())
	str += fmt.Sprintf(" - Pause (ms.): %d\n", c.Pause)
	str += fmt.Sprintf(" - Sample Rate: %d\n", c.SampleRate)
	str += fmt.Sprintf(" - Compression: %d\n", c.CompressionType)
	str += fmt.Sprintf(" - Pulse Count: %d\n", c.StoredPulseCount)

	return str
}
