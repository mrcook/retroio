package tzx

import (
	"fmt"

	"github.com/mrcook/tzxbrowser/tape"
)

// JumpTo
// ID: 23h (35d)
// This block will enable you to jump from one block to another within the file. The value is a
// signed short word (usually 'signed short' in C); Some examples:
//   - Jump 0 = 'Loop Forever' - this should never happen
//   - Jump 1 = 'Go to the next block' - it is like NOP in assembler ;)
//   - Jump 2 = 'Skip one block'
//   - Jump -1 = 'Go to the previous block'
// All blocks are included in the block count!
type JumpTo struct {
	Value int16 // WORD  Relative jump value
}

func (j *JumpTo) Process(file *tape.File) {
	j.Value = file.ReadSignedShort()
}

func (j JumpTo) Id() int {
	return 35
}

func (j JumpTo) Name() string {
	return "Jump To"
}

// Metadata returns a human readable string of the block data
func (j JumpTo) Metadata() string {
	return fmt.Sprintf("> %-19s : %d", j.Name(), j.Value)
}
