package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// Message
// ID: 31h (49d)
// This will enable the emulators to display a message for a given time. This should not stop the
// tape and it should not make silence. If the time is 0 then the emulator should wait for the user
// to press a key.
// The text message should:
//   - stick to a maximum of 30 chars per line;
//   - use single 0x0D (13 decimal) to separate lines;
//   - stick to a maximum of 8 lines.
type Message struct {
	DisplayTime uint8  // BYTE  Time (in seconds) for which the message should be displayed
	Length      uint8  // N BYTE  Length of the text message
	Message     []byte // CHAR[N] Message that should be displayed in ASCII format
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (m *Message) Read(file *tape.File) {
	m.DisplayTime, _ = file.ReadByte()
	m.Length, _ = file.ReadByte()

	for _, b := range file.ReadBytes(int(m.Length)) {
		m.Message = append(m.Message, b)
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (m Message) Id() uint8 {
	return 0x31
}

// Name of the block as given in the TZX specification.
func (m Message) Name() string {
	return "Message"
}

// ToString returns a human readable string of the block data
func (m Message) ToString() string {
	str := fmt.Sprintf("> %-19s : display for %d seconds\n", m.Name(), m.DisplayTime)
	str += fmt.Sprintf(" - Message: %s\n", m.Message)
	return str
}
