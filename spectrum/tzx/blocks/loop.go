package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/storage"
)

// LoopStart
// ID: 24h (36d)
// If you have a sequence of identical blocks, or of identical groups of blocks, you can use
// this block to tell how many times they should be repeated. This block is the same as the
// FOR statement in BASIC.
type LoopStart struct {
	RepetitionCount uint16 // WORD  Number of repetitions (greater than 1)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (l *LoopStart) Read(reader *storage.Reader) {
	l.RepetitionCount = reader.ReadShort()
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (l LoopStart) Id() uint8 {
	return 0x24
}

// Name of the block as given in the TZX specification.
func (l LoopStart) Name() string {
	return "Loop Start"
}

func (l LoopStart) BlockData() tap.BlockI {
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
type LoopEnd struct{}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (l *LoopEnd) Read(reader *storage.Reader) {}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (l LoopEnd) Id() uint8 {
	return 0x25
}

// Name of the block as given in the TZX specification.
func (l LoopEnd) Name() string {
	return "Loop End"
}

func (l LoopEnd) BlockData() tap.BlockI {
	return nil
}

// String returns a human readable string of the block data
func (l LoopEnd) String() string {
	return fmt.Sprintf("%s", l.Name())
}
