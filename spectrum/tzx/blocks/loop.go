package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
)

// LoopStart
// ID: 24h (36d)
// If you have a sequence of identical blocks, or of identical groups of blocks, you can use
// this block to tell how many times they should be repeated. This block is the same as the
// FOR statement in BASIC.
type LoopStart struct {
	BlockID         types.BlockType
	RepetitionCount uint16 // Number of repetitions (greater than 1)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (l *LoopStart) Read(reader *storage.Reader) error {
	l.BlockID = types.BlockType(reader.ReadByte())
	if l.BlockID != l.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", l.Id(), l.BlockID)
	}

	l.RepetitionCount = reader.ReadShort()

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (l LoopStart) Id() types.BlockType {
	return types.LoopStart
}

// Name of the block as given in the TZX specification.
func (l LoopStart) Name() string {
	return "Loop Start"
}

func (l LoopStart) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (l LoopStart) String() string {
	return fmt.Sprintf("%-19s : %d times", l.Name(), l.RepetitionCount)
}

// LoopEnd
// ID: 25h (37d)
// This is the same as BASIC's NEXT statement. It means that the utility should jump back to the
// start of the loop if it hasn't been run for the specified number of times.
// This block has no body.
type LoopEnd struct {
	BlockID types.BlockType
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (l *LoopEnd) Read(reader *storage.Reader) error {
	l.BlockID = types.BlockType(reader.ReadByte())
	if l.BlockID != l.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", l.Id(), l.BlockID)
	}
	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (l LoopEnd) Id() types.BlockType {
	return types.LoopEnd
}

// Name of the block as given in the TZX specification.
func (l LoopEnd) Name() string {
	return "Loop End"
}

func (l LoopEnd) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (l LoopEnd) String() string {
	return fmt.Sprintf("%s", l.Name())
}
