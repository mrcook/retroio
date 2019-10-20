package blocks

import (
	"bufio"
	"fmt"

	"mrcook/tzxit/tape"
)

// DirectRecording
// ID: 15h (21d)
// This block is used for tapes which have some parts in a format such that the turbo loader block
// cannot be used. This is not like a VOC file, since the information is much more compact. Each
// sample value is represented by one bit only (0 for low, 1 for high) which means that the block
// will be at most 1/8 the size of the equivalent VOC.
// The preferred sampling frequencies are 22050 or 44100 Hz (158 or 79 T-states/sample). Please,
// if you can, don't use other sampling frequencies.
// Please use this block only if you cannot use any other block.
type DirectRecording struct {
	TStatesPerSample uint16  // WORD      Number of T-states per sample (bit of data)
	Pause            uint16  // WORD      Pause after this block in milliseconds (ms.)
	UsedBits         uint8   // BYTE      Used bits (samples) in last byte of data (1-8) (e.g. if this is 2, only first two samples of the last byte will be played)
	Length           uint32  // N BYTE[3] Length of samples' data
	Data             []uint8 // BYTE[N]   Samples data. Each bit represents a state on the EAR port (i.e. one sample). MSb is played first.
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (d *DirectRecording) Read(reader *bufio.Reader) {
	d.TStatesPerSample = tape.ReadShort(reader)
	d.Pause = tape.ReadShort(reader)
	d.UsedBits, _ = reader.ReadByte()

	length := tape.ReadNextBytes(reader, 3)
	length = append(length, 0) // add 4th byte
	d.Length = tape.BytesToLong(length)

	// Yep, we're discarding the data for the moment
	tape.ReadNextBytes(reader, int(d.Length))
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (d DirectRecording) Id() uint8 {
	return 0x15
}

// Name of the block as given in the TZX specification.
func (d DirectRecording) Name() string {
	return "Direct Recording"
}

// ToString returns a human readable string of the block data
func (d DirectRecording) ToString() string {
	return fmt.Sprintf("> %-19s : %d T-States, %d bytes", d.Name(), d.TStatesPerSample, d.Length)
}
