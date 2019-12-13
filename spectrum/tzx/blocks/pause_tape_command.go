package blocks

import (
	"fmt"

	"retroio/storage"
)

// PauseTapeCommand
// ID: 20h (32d)
// This will make a silence (low amplitude level (0)) for a given time in milliseconds. If the
// value is 0 then the emulator or utility should (in effect) STOP THE TAPE, i.e. should not
// continue loading until the user or emulator requests it.
type PauseTapeCommand struct {
	Pause uint16 // WORD  Pause duration (ms.)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PauseTapeCommand) Read(reader *storage.Reader) {
	p.Pause = reader.ReadShort()
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PauseTapeCommand) Id() uint8 {
	return 0x20
}

// Name of the block as given in the TZX specification.
func (p PauseTapeCommand) Name() string {
	return "Pause Tape Command"
}

// String returns a human readable string of the block data
func (p PauseTapeCommand) String() string {
	return fmt.Sprintf("%-19s : %d ms.", p.Name(), p.Pause)
}
