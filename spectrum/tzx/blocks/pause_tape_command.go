package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// PauseTapeCommand
// ID: 20h (32d)
// This will make a silence (low amplitude level (0)) for a given time in milliseconds. If the
// value is 0 then the emulator or utility should (in effect) STOP THE TAPE, i.e. should not
// continue loading until the user or emulator requests it.
type PauseTapeCommand struct {
	BlockID types.BlockType
	Pause   uint16 // Pause duration (ms.)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PauseTapeCommand) Read(reader *storage.Reader) error {
	p.BlockID = types.BlockType(reader.ReadByte())
	if p.BlockID != p.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", p.Id(), p.BlockID)
	}

	p.Pause = reader.ReadShort()

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PauseTapeCommand) Id() types.BlockType {
	return types.PauseTapeCommand
}

// Name of the block as given in the TZX specification.
func (p PauseTapeCommand) Name() string {
	return "Pause Tape Command"
}

func (p PauseTapeCommand) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (p PauseTapeCommand) String() string {
	return fmt.Sprintf("%-19s : %d ms.", p.Name(), p.Pause)
}
