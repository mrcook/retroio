package csw2

import (
	"fmt"

	"github.com/mrcook/retroio/librio/storage"
)

// Length of this block size, without the CSW2 data itself.
// This is required to know how much CSW data to read.
const CswHeaderSize uint32 = 52

type CSW struct {
	// 23-byte signature: `Compressed Square Wave` + 0x1A terminator byte
	Signature [23]byte

	// New instances should always be CSW-2 (2.0)
	MajorVersion uint8
	MinorVersion uint8

	SampleRate       uint32   // Sample rate
	StoredPulseCount uint32   // Total number of pulses (after decompression)
	CompressionType  uint8    // Compression type: 0x01: RLE, 0x02: Z-RLE
	Flags            uint8    // Flags. b0: initial polarity; if set, the signal starts at logical high
	HdrLength        uint8    // Header extension length in bytes (always 0x00).
	Description      [16]byte // Encoding application description. Information about the tool which created the file (e.g. name and version)

	// NOTE: the specification states the HDR is "For future expansions only".
	// As v2 does not support this, it will always be "not present".
	// hdrData []uint8 // Header extension data (if present)

	// The actual CSW data
	// NOTE: this is intentionally stored as plain `byte` data - this may change.
	Data []byte
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CSW) Read(reader *storage.Reader, cswBlockSize int) error {
	dataLength := cswBlockSize - int(CswHeaderSize)
	if dataLength < 0 {
		return fmt.Errorf("invalid CSW2 block length")
	}

	for i, b := range reader.ReadBytes(23) {
		c.Signature[i] = b
	}

	c.MajorVersion = reader.ReadByte()
	c.MinorVersion = reader.ReadByte()
	c.SampleRate = reader.ReadLong()
	c.StoredPulseCount = reader.ReadLong()
	c.CompressionType = reader.ReadByte()
	c.Flags = reader.ReadByte()
	c.HdrLength = reader.ReadByte()

	for i, b := range reader.ReadBytes(16) {
		c.Description[i] = b
	}

	c.Data = make([]byte, dataLength)
	if _, err := reader.Read(c.Data); err != nil {
		return err
	}

	return nil
}
