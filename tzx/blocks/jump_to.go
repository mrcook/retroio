package blocks

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
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

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (j *JumpTo) Read(file *tape.Reader) {
	j.Value = file.ReadSignedShort()
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (j JumpTo) Id() uint8 {
	return 0x23
}

// Name of the block as given in the TZX specification.
func (j JumpTo) Name() string {
	return "Jump To"
}

// ToString returns a human readable string of the block data
func (j JumpTo) ToString() string {
	return fmt.Sprintf("> %-19s : %d", j.Name(), j.Value)
}
