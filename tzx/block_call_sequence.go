package tzx

import (
	"fmt"

	"github.com/mrcook/tzxbrowser/tape"
)

// CallSequence
// ID: 26h (38d)
// This block is an analogue of the CALL Subroutine statement. It basically executes a sequence of
// blocks that are somewhere else and then goes back to the next block. Because more than one call
// can be normally used you can include a list of sequences to be called. The 'nesting' of call
// blocks is also not allowed for the simplicity reasons. You can, of course, use the CALL blocks
// in the LOOP sequences and vice versa. The value is relative for the obvious reasons - so that
// you can add some blocks in the beginning of the file without disturbing the call values. Please
// take a look at 'Jump To Block' for reference on the values.
type CallSequence struct {
	Count  uint16   // N WORD  Number of calls to be made
	Blocks []uint16 // WORD[N] Array of call block numbers (relative-signed offsets)
}

func (c *CallSequence) Process(file *tape.File) {
	c.Count = file.ReadShort()

	for i := 0; i < int(c.Count); i++ {
		c.Blocks = append(c.Blocks, file.ReadShort())
	}
}

func (c CallSequence) Id() int {
	return 38
}

func (c CallSequence) Name() string {
	return "Call Sequence"
}

// ToString returns a human readable string of the block data
func (c CallSequence) ToString() string {
	str := fmt.Sprintf("> %s\n", c.Name())
	for _, b := range c.Blocks {
		str += fmt.Sprintf(" - %d\n", b)
	}
	return str
}

// ReturnFromSequence
// ID: 27h (39d)
// This block indicates the end of the Called Sequence. The next block played will be the block after
// the last CALL block (or the next Call, if the Call block had multiple calls).
// This block has no body.
type ReturnFromSequence struct{}

func (r ReturnFromSequence) Process(file *tape.File) {}

func (r ReturnFromSequence) Id() int {
	return 39
}

func (r ReturnFromSequence) Name() string {
	return "Return from Sequence"
}

// ToString returns a human readable string of the block data
func (r ReturnFromSequence) ToString() string {
	return r.Name()
}
