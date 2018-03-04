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

func (m *CustomInfo) Process(file *File) {
	for i, b := range file.ReadBytes(10) {
		m.Identification[i] = b
	}

	m.Length = file.ReadLong()

	for _, b := range file.ReadBytes(int(m.Length)) {
		m.Info = append(m.Info, b)
	}
}

func (m CustomInfo) Id() int {
	return 53
}

func (m CustomInfo) Name() string {
	return "Custom Info"
}

// Metadata returns a human readable string of the block data
func (m CustomInfo) Metadata() string {
	str := ""
	str += fmt.Sprintf("Identification: %s\n", m.Identification)
	str += fmt.Sprintf("Length:         %d\n", m.Length)
	str += fmt.Sprintf("Info:           %s\n", m.Info)
	return str
}
