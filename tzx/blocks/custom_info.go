package blocks

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// CustomInfo
// ID: 35h (53d)
// This block can be used to save any information you want. For example, it might contain some
// information written by a utility, extra settings required by a particular emulator, or even
// poke data.
type CustomInfo struct {
	Identification [10]byte // CHAR[10]  Identification string (in ASCII)
	Length         uint32   // L DWORD   Length of the custom info
	Info           []uint8  // BYTE[L]   Custom info
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (c *CustomInfo) Read(file *tape.Reader) {
	for i, b := range file.ReadBytes(10) {
		c.Identification[i] = b
	}

	c.Length = file.ReadLong()

	for _, b := range file.ReadBytes(int(c.Length)) {
		c.Info = append(c.Info, b)
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (c CustomInfo) Id() uint8 {
	return 0x35
}

// Name of the block as given in the TZX specification.
func (c CustomInfo) Name() string {
	return "Custom Info"
}

// ToString returns a human readable string of the block data
func (c CustomInfo) ToString() string {
	return fmt.Sprintf("> %-19s : %s - %s", c.Name(), c.Identification, c.Info)
}
