package blocks

import (
	"bufio"
	"fmt"

	"mrcook/tzxit/tap"
	"mrcook/tzxit/tape"
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
	DataBlock tap.DataBlock
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (t *TurboSpeedData) Read(reader *bufio.Reader) {
	t.PilotPulse = tape.ReadShort(reader)
	t.SyncFirstPulse = tape.ReadShort(reader)
	t.SyncSecondPulse = tape.ReadShort(reader)
	t.ZeroBitPulse = tape.ReadShort(reader)
	t.OneBitPulse = tape.ReadShort(reader)
	t.PilotTone = tape.ReadShort(reader)
	t.UsedBits, _ = reader.ReadByte()
	t.Pause = tape.ReadShort(reader)

	length := tape.ReadNextBytes(reader, 3)
	length = append(length, 0) // add 4th byte
	t.Length = tape.BytesToLong(length)

	// Yep, we're discarding the data for the moment
	tape.ReadNextBytes(reader, int(t.Length))
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (t TurboSpeedData) Id() uint8 {
	return 0x11
}

// Name of the block as given in the TZX specification.
func (t TurboSpeedData) Name() string {
	return "Turbo Speed Data"
}

// ToString returns a human readable string of the block data
func (t TurboSpeedData) ToString() string {
	return fmt.Sprintf("> %-19s : %d bytes, pause for %d ms.", t.Name(), t.Length, t.Pause)
}
