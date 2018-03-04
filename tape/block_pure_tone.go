package tape

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

// PureTone
// ID: 12h (18d)
// This will produce a tone which is basically the same as the pilot tone in the ID 10, ID 11
// blocks. You can define how long the pulse is and how many pulses are in the tone.
type PureTone struct {
	Length     uint16 // WORD Length of one pulse in T-states
	PulseCount uint16 // WORD Number of pulses
}

func (t *PureTone) Process(file *File) {
	data := file.ReadBytes(4)

	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, &t)
	if err != nil {
		log.Fatal("binary.Read failed: ", err)
	}
}

func (t PureTone) Id() int {
	return 18
}

func (t PureTone) Name() string {
	return "Pure Tone"
}

// Metadata returns a human readable string of the block data
func (t PureTone) Metadata() string {
	return fmt.Sprintf("Length:     %d\nPulseCount: %d\n", t.Length, t.PulseCount)
}
