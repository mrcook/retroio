package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/storage"
)

// PureTone
// ID: 12h (18d)
// This will produce a tone which is basically the same as the pilot tone in the ID 10, ID 11
// blocks. You can define how long the pulse is and how many pulses are in the tone.
type PureTone struct {
	Length     uint16 // WORD Length of one pulse in T-states
	PulseCount uint16 // WORD Number of pulses
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PureTone) Read(reader *storage.Reader) {
	p.Length = reader.ReadShort()
	p.PulseCount = reader.ReadShort()
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PureTone) Id() uint8 {
	return 0x12
}

// Name of the block as given in the TZX specification.
func (p PureTone) Name() string {
	return "Pure Tone"
}

func (p PureTone) BlockData() tap.BlockI {
	return nil
}

// String returns a human readable string of the block data
func (p PureTone) String() string {
	return fmt.Sprintf("%-19s : %d T-States, %d pulses", p.Name(), p.Length, p.PulseCount)
}
