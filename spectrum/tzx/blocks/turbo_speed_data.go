package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
)

// TurboSpeedData
// ID: 11h (17d)
// This block is very similar to the normal TAP block but with some additional info on the timings
// and other important differences. The same tape encoding is used as for the standard speed data
// block. If a block should use some non-standard sync or pilot tones (i.e. all sorts of protection
// schemes) then use the next three blocks to describe it.
type TurboSpeedData struct {
	BlockID         types.BlockType
	PilotPulse      uint16   // Length of PILOT pulse {2168}
	SyncFirstPulse  uint16   // Length of SYNC first pulse {667}
	SyncSecondPulse uint16   // Length of SYNC second pulse {735}
	ZeroBitPulse    uint16   // Length of ZERO bit pulse {855}
	OneBitPulse     uint16   // Length of ONE bit pulse {1710}
	PilotTone       uint16   // Length of PILOT tone (number of pulses) {8063 header (flag<128), 3223 data (flag>=128)}
	UsedBits        uint8    // Used bits in the last byte (other bits should be 0) {8} (e.g. if this is 6, then the bits used (x) in the last byte are: xxxxxx00, where MSb is the leftmost bit, LSb is the rightmost bit)
	Pause           uint16   // Pause after this block (ms.) {1000}
	Length          [3]uint8 // Length of data that follows.
	DataBlock       []uint8  // Data as in .TAP files

	displayLength uint32
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (t *TurboSpeedData) Read(reader *storage.Reader) error {
	t.BlockID = types.BlockType(reader.ReadByte())
	if t.BlockID != t.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", t.Id(), t.BlockID)
	}

	t.PilotPulse = reader.ReadShort()
	t.SyncFirstPulse = reader.ReadShort()
	t.SyncSecondPulse = reader.ReadShort()
	t.ZeroBitPulse = reader.ReadShort()
	t.OneBitPulse = reader.ReadShort()
	t.PilotTone = reader.ReadShort()
	t.UsedBits = reader.ReadByte()
	t.Pause = reader.ReadShort()

	copy(t.Length[:], reader.ReadBytes(3))

	t.displayLength = reader.Bytes3ToLong(t.Length)

	// TODO: read this as TAP data.
	t.DataBlock = make([]byte, t.displayLength)
	_, err := reader.Read(t.DataBlock)
	return err
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (t TurboSpeedData) Id() types.BlockType {
	return types.TurboSpeedData
}

// Name of the block as given in the TZX specification.
func (t TurboSpeedData) Name() string {
	return "Turbo Speed Data"
}

func (t TurboSpeedData) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (t TurboSpeedData) String() string {
	return fmt.Sprintf("%-19s : %d bytes, pause for %d ms.", t.Name(), t.displayLength, t.Pause)
}
