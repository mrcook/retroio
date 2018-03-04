package tape

import "fmt"

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

func (m *Message) Process(file *File) {
	m.DisplayTime, _ = file.ReadByte()
	m.Length, _ = file.ReadByte()

	for _, b := range file.ReadBytes(int(m.Length)) {
		m.Message = append(m.Message, b)
	}
}

func (m Message) Id() int {
	return 49
}

func (m Message) Name() string {
	return "Message"
}

// Metadata returns a human readable string of the block data
func (m Message) Metadata() string {
	str := fmt.Sprintf("> %s\n", m.Name())
	str += fmt.Sprintf(" - Display for %d seconds\n", m.DisplayTime)
	str += fmt.Sprintf(" - Message: %s\n", m.Message)
	return str
}
