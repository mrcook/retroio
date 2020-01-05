package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
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
	BlockID types.BlockType
	Value   int16 // Relative jump value
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (j *JumpTo) Read(reader *storage.Reader) error {
	j.BlockID = types.BlockType(reader.ReadByte())
	if j.BlockID != j.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", j.Id(), j.BlockID)
	}

	j.Value = int16(reader.ReadShort())

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (j JumpTo) Id() types.BlockType {
	return types.JumpTo
}

// Name of the block as given in the TZX specification.
func (j JumpTo) Name() string {
	return "Jump To"
}

func (j JumpTo) BlockData() tap.BlockI {
	return nil
}

// String returns a human readable string of the block data
func (j JumpTo) String() string {
	return fmt.Sprintf("%-19s : %d", j.Name(), j.Value)
}
