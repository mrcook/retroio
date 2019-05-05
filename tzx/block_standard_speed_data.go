package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tap"
	"github.com/mrcook/tzxit/tape"
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

func (s *StandardSpeedData) Read(file *tape.File) {
	s.Pause = file.ReadShort()
	s.Length = file.ReadShort()

	// For the moment, discard data unless it's a TAP block
	data := file.ReadBytes(int(s.Length))
	if s.Length == 19 {
		s.Data = data
	}
}

func (s StandardSpeedData) Id() int {
	return 16
}

func (s StandardSpeedData) Name() string {
	return "Standard Speed Data"
}

// ToString returns a human readable string of the block data
func (s StandardSpeedData) ToString() string {
	str := fmt.Sprintf("> %-19s : %d bytes, pause for %d ms.", s.Name(), s.Length, s.Pause)

	if s.Length == 19 {
		b, err := tap.Unmarshal(s.Data)
		if err != nil {
			str += fmt.Sprintf("TAP BLOCK ERROR: %v\n", err)
		} else {
			str += "\n" + b.ToString()
		}
	}

	return str
}
