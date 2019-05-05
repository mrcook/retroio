package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// PureTone
// ID: 12h (18d)
// This will produce a tone which is basically the same as the pilot tone in the ID 10, ID 11
// blocks. You can define how long the pulse is and how many pulses are in the tone.
type PureTone struct {
	Length     uint16 // WORD Length of one pulse in T-states
	PulseCount uint16 // WORD Number of pulses
}

func (p *PureTone) Read(file *tape.File) {
	p.Length = file.ReadShort()
	p.PulseCount = file.ReadShort()
}

func (p PureTone) Id() int {
	return 18
}

func (p PureTone) Name() string {
	return "Pure Tone"
}

// ToString returns a human readable string of the block data
func (p PureTone) ToString() string {
	return fmt.Sprintf("> %-19s : %d T-States, %d pulses", p.Name(), p.Length, p.PulseCount)
}
