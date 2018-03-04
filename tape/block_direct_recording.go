package tape

import "fmt"

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
	Length           uint16  // N BYTE[3] Length of samples' data
	LengthSpareByte  uint8   // NOTE: `length` above uses only 2-bytes for value but specification says 3-bytes, so this is for the spare.
	Data             []uint8 // BYTE[N]   Samples data. Each bit represents a state on the EAR port (i.e. one sample). MSb is played first.
}

func (d *DirectRecording) Process(file *File) {
	d.TStatesPerSample = file.ReadShort()
	d.Pause = file.ReadShort()
	d.UsedBits, _ = file.ReadByte()
	d.Length = file.ReadShort()
	d.LengthSpareByte, _ = file.ReadByte()

	// Yep, we're discarding the data for the moment
	file.ReadBytes(int(d.Length))
}

func (d DirectRecording) Id() int {
	return 21
}

func (d DirectRecording) Name() string {
	return "Direct Recording"
}

// Metadata returns a human readable string of the block data
func (d DirectRecording) Metadata() string {
	str := fmt.Sprintf("> %s\n", d.Name())
	str += fmt.Sprintf(" - T-States per Sample: %d\n", d.TStatesPerSample)
	str += fmt.Sprintf(" - Used Bits:           %d\n", d.UsedBits)
	str += fmt.Sprintf(" - Pause (ms.):         %d\n", d.Pause)
	return str
}
