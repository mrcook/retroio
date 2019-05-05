package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// LoopStart
// ID: 24h (36d)
// If you have a sequence of identical blocks, or of identical groups of blocks, you can use
// this block to tell how many times they should be repeated. This block is the same as the
// FOR statement in BASIC.
type LoopStart struct {
	RepetitionCount uint16 // WORD  Number of repetitions (greater than 1)
}

func (l *LoopStart) Read(file *tape.File) {
	l.RepetitionCount = file.ReadShort()
}

func (l LoopStart) Id() uint8 {
	return 0x24
}

func (l LoopStart) Name() string {
	return "Loop Start"
}

// ToString returns a human readable string of the block data
func (l LoopStart) ToString() string {
	return fmt.Sprintf("> %-19s : %d times", l.Name(), l.RepetitionCount)
}

// LoopEnd
// ID: 25h (37d)
// This is the same as BASIC's NEXT statement. It means that the utility should jump back to the
// start of the loop if it hasn't been run for the specified number of times.
// This block has no body.
type LoopEnd struct{}

func (l *LoopEnd) Read(file *tape.File) {}

func (l LoopEnd) Id() uint8 {
	return 0x25
}

func (l LoopEnd) Name() string {
	return "Loop End"
}

// ToString returns a human readable string of the block data
func (l LoopEnd) ToString() string {
	return fmt.Sprintf("> %s", l.Name())
}
