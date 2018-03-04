package tape

import (
	"fmt"
)

// PauseTapeCommand
// ID: 20h (32d)
// This will make a silence (low amplitude level (0)) for a given time in milliseconds. If the
// value is 0 then the emulator or utility should (in effect) STOP THE TAPE, i.e. should not
// continue loading until the user or emulator requests it.
type PauseTapeCommand struct {
	Pause uint16 // WORD  Pause duration (ms.)
}

func (p *PauseTapeCommand) Process(file *File) {
	p.Pause = file.ReadShort()
}

func (p PauseTapeCommand) Id() int {
	return 32
}

func (p PauseTapeCommand) Name() string {
	return "Pause Tape Command"
}

// Metadata returns a human readable string of the block data
func (p PauseTapeCommand) Metadata() string {
	return fmt.Sprintf("> %s : %d ms.\n", p.Name(), p.Pause)
}
