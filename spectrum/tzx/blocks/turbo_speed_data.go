package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/storage"
)

// TurboSpeedData
// ID: 11h (17d)
// This block is very similar to the normal TAP block but with some additional info on the timings
// and other important differences. The same tape encoding is used as for the standard speed data
// block. If a block should use some non-standard sync or pilot tones (i.e. all sorts of protection
// schemes) then use the next three blocks to describe it.
type TurboSpeedData struct {
	PilotPulse      uint16 // WORD      Length of PILOT pulse {2168}
	SyncFirstPulse  uint16 // WORD      Length of SYNC first pulse {667}
	SyncSecondPulse uint16 // WORD      Length of SYNC second pulse {735}
	ZeroBitPulse    uint16 // WORD      Length of ZERO bit pulse {855}
	OneBitPulse     uint16 // WORD      Length of ONE bit pulse {1710}
	PilotTone       uint16 // WORD      Length of PILOT tone (number of pulses) {8063 header (flag<128), 3223 data (flag>=128)}
	UsedBits        uint8  // BYTE      Used bits in the last byte (other bits should be 0) {8} (e.g. if this is 6, then the bits used (x) in the last byte are: xxxxxx00, where MSb is the leftmost bit, LSb is the rightmost bit)
	Pause           uint16 // WORD      Pause after this block (ms.) {1000}

	Length uint32 // N BYTE[3] Length of data that follows. NOTE the use of a DWORD for the property type

	// A single .TAP DataBlock consisting of:
	//   WORD    Length of data that follows
	//   BYTE[N] Data as in .TAP files
	DataBlock tap.BlockI
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (t *TurboSpeedData) Read(reader *storage.Reader) {
	t.PilotPulse = reader.ReadShort()
	t.SyncFirstPulse = reader.ReadShort()
	t.SyncSecondPulse = reader.ReadShort()
	t.ZeroBitPulse = reader.ReadShort()
	t.OneBitPulse = reader.ReadShort()
	t.PilotTone = reader.ReadShort()
	t.UsedBits = reader.ReadByte()
	t.Pause = reader.ReadShort()

	length := reader.ReadNextBytes(3)
	length = append(length, 0) // add 4th byte
	t.Length = reader.BytesToLong(length)

	// Yep, we're discarding the data for the moment
	data := make([]byte, t.Length)
	_, _ = reader.Read(data)
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (t TurboSpeedData) Id() uint8 {
	return 0x11
}

// Name of the block as given in the TZX specification.
func (t TurboSpeedData) Name() string {
	return "Turbo Speed Data"
}

func (t TurboSpeedData) BlockData() tap.BlockI {
	return t.DataBlock
}

// String returns a human readable string of the block data
func (t TurboSpeedData) String() string {
	return fmt.Sprintf("%-19s : %d bytes, pause for %d ms.", t.Name(), t.Length, t.Pause)
}
