package tzx

import (
	"fmt"

	"github.com/mrcook/tzxbrowser/tape"
)

// LoopStart
// ID: 24h (36d)
// If you have a sequence of identical blocks, or of identical groups of blocks, you can use
// this block to tell how many times they should be repeated. This block is the same as the
// FOR statement in BASIC.
type LoopStart struct {
	RepetitionCount uint16 // WORD  Number of repetitions (greater than 1)
}

func (l *LoopStart) Process(file *tape.File) {
	l.RepetitionCount = file.ReadShort()
}

func (l LoopStart) Id() int {
	return 36
}

func (l LoopStart) Name() string {
	return "Loop Start"
}

// Metadata returns a human readable string of the block data
func (l LoopStart) Metadata() string {
	return fmt.Sprintf("> %-19s : %d times", l.Name(), l.RepetitionCount)
}

// LoopEnd
// ID: 25h (37d)
// This is the same as BASIC's NEXT statement. It means that the utility should jump back to the
// start of the loop if it hasn't been run for the specified number of times.
// This block has no body.
type LoopEnd struct{}

func (l *LoopEnd) Process(file *tape.File) {}

func (l LoopEnd) Id() int {
	return 37
}

func (l LoopEnd) Name() string {
	return "Loop End"
}

// Metadata returns a human readable string of the block data
func (l LoopEnd) Metadata() string {
	return fmt.Sprintf("> %s", l.Name())
}
