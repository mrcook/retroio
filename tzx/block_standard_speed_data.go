package tzx

import (
	"fmt"

	"github.com/mrcook/tzxbrowser/tape"
)

// StandardSpeedData
// ID: 10h (16d)
// This block must be replayed with the standard Spectrum ROM timing values - see the values in
// curly brackets in block ID 11. The pilot tone consists in 8063 pulses if the first data byte
// (flag byte) is < 128, 3223 otherwise. This block can be used for the ROM loading routines AND
// for custom loading routines that use the same timings as ROM ones do.
type StandardSpeedData struct {
	Pause  uint16  // WORD    Pause after this block (ms.) {1000}
	Length uint16  // WORD    Length of data that follow
	Data   []uint8 // BYTE[N] Data as in .TAP files
}

func (s *StandardSpeedData) Process(file *tape.File) {
	s.Pause = file.ReadShort()
	s.Length = file.ReadShort()

	// Yep, we're discarding the data for the moment
	file.ReadBytes(int(s.Length))
}

func (s StandardSpeedData) Id() int {
	return 16
}

func (s StandardSpeedData) Name() string {
	return "Standard Speed Data"
}

// Metadata returns a human readable string of the block data
func (s StandardSpeedData) Metadata() string {
	return fmt.Sprintf("> %-19s : %d bytes, pause for %d ms.", s.Name(), s.Length, s.Pause)
}
