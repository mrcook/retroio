package blocks

import (
	"bufio"
	"fmt"

	"retroio/tap"
	"retroio/tape"
)

// StandardSpeedData
// ID: 10h (16d)
// This block must be replayed with the standard Spectrum ROM timing values - see the values in
// curly brackets in block ID 11. The pilot tone consists in 8063 pulses if the first data byte
// (flag byte) is < 128, 3223 otherwise. This block can be used for the ROM loading routines AND
// for custom loading routines that use the same timings as ROM ones do.
type StandardSpeedData struct {
	Pause uint16 // WORD    Pause after this block (ms.) {1000}

	// A single .TAP DataBlock consisting of:
	//   WORD    Length of data that follows
	//   BYTE[N] Data as in .TAP files
	DataBlock tap.DataBlock
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *StandardSpeedData) Read(reader *bufio.Reader) {
	s.Pause = tape.ReadShort(reader)

	t := tap.NewReader(reader)
	length, _ := tape.PeekBlockLength(reader)
	if length == 19 {
		s.DataBlock, _ = t.ReadHeaderBlock()
	} else {
		s.DataBlock, _ = t.ReadDataBlock()
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s StandardSpeedData) Id() uint8 {
	return 0x10
}

// Name of the block as given in the TZX specification.
func (s StandardSpeedData) Name() string {
	return "Standard Speed Data"
}

// ToString returns a human readable string of the block data
func (s StandardSpeedData) ToString() string {
	str := fmt.Sprintf("> %-19s : %d bytes, pause for %d ms.", s.Name(), s.DataBlock.Length, s.Pause)
	str += fmt.Sprintf("\n    %s", s.DataBlock.TapeData.ToString())

	return str
}
