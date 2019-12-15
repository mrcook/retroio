package blocks

import (
	"fmt"
	"io"

	"retroio/spectrum/tap"
	"retroio/storage"
)

// PureData
// ID: 14h (20d)
// This is the same as in the turbo loading data block, except that it has no pilot or sync pulses.
type PureData struct {
	ZeroBitPulse uint16 // WORD      Length of ZERO bit pulse
	OneBitPulse  uint16 // WORD      Length of ONE bit pulse
	UsedBits     uint8  // BYTE      Used bits in last byte (other bits should be 0) (e.g. if this is 6, then the bits used (x) in the last byte are: xxxxxx00, where MSb is the leftmost bit, LSb is the rightmost bit)
	Pause        uint16 // WORD      Pause after this block (ms.)

	Length uint32 // N BYTE[3] Length of data that follows. NOTE the use of a DWORD for the property type

	// A single .TAP DataBlock consisting of:
	//   WORD    Length of data that follows
	//   BYTE[N] Data as in .TAP files
	DataBlock tap.BlockI
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PureData) Read(reader *storage.Reader) {
	p.OneBitPulse = reader.ReadShort()
	p.ZeroBitPulse = reader.ReadShort()
	p.UsedBits = reader.ReadByte()
	p.Pause = reader.ReadShort()

	length := reader.ReadNextBytes(3)
	length = append(length, 0) // add 4th byte
	p.Length = reader.BytesToLong(length)

	// Yep, we're discarding the data for the moment
	data := make([]byte, p.Length)
	_, _ = io.ReadFull(reader, data)
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PureData) Id() uint8 {
	return 0x14
}

// Name of the block as given in the TZX specification.
func (p PureData) Name() string {
	return "Pure Data"
}

func (p PureData) BlockData() tap.BlockI {
	return p.DataBlock
}

// String returns a human readable string of the block data
func (p PureData) String() string {
	return fmt.Sprintf("%-19s : %d bytes, pause for %d ms.", p.Name(), p.Length, p.Pause)
}
