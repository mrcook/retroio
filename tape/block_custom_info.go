package tape

import "fmt"

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

func (c *CustomInfo) Process(file *File) {
	for i, b := range file.ReadBytes(10) {
		c.Identification[i] = b
	}

	c.Length = file.ReadLong()

	for _, b := range file.ReadBytes(int(c.Length)) {
		c.Info = append(c.Info, b)
	}
}

func (c CustomInfo) Id() int {
	return 53
}

func (c CustomInfo) Name() string {
	return "Custom Info"
}

// Metadata returns a human readable string of the block data
func (c CustomInfo) Metadata() string {
	str := fmt.Sprintf("> %s\n", c.Name())
	str += fmt.Sprintf(" - Type: %s\n", c.Identification)
	str += fmt.Sprintf(" - Info: %s\n", c.Info)
	return str
}
