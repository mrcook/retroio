package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/storage"
)

// PureData
// ID: 14h (20d)
// This is the same as in the turbo loading data block, except that it has no pilot or sync pulses.
type PureData struct {
	BlockID      types.BlockType
	ZeroBitPulse uint16   // Length of ZERO bit pulse
	OneBitPulse  uint16   // Length of ONE bit pulse
	UsedBits     uint8    // Used bits in last byte (other bits should be 0) (e.g. if this is 6, then the bits used (x) in the last byte are: xxxxxx00, where MSb is the leftmost bit, LSb is the rightmost bit)
	Pause        uint16   // Pause after this block (ms.)
	Length       [3]uint8 // Length of data that follows.
	DataBlock    []uint8  // Data as in .TAP files

	displayLength uint32
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (p *PureData) Read(reader *storage.Reader) error {
	p.BlockID = types.BlockType(reader.ReadByte())
	if p.BlockID != p.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", p.Id(), p.BlockID)
	}

	p.OneBitPulse = reader.ReadShort()
	p.ZeroBitPulse = reader.ReadShort()
	p.UsedBits = reader.ReadByte()
	p.Pause = reader.ReadShort()
	copy(p.Length[:], reader.ReadBytes(3))

	p.displayLength = reader.Bytes3ToLong(p.Length)

	// TODO: read this as TAP data.
	p.DataBlock = make([]byte, p.displayLength)
	_, err := reader.Read(p.DataBlock)
	return err
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (p PureData) Id() types.BlockType {
	return types.PureData
}

// Name of the block as given in the TZX specification.
func (p PureData) Name() string {
	return "Pure Data"
}

func (p PureData) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
func (p PureData) String() string {
	return fmt.Sprintf("%-19s : %d bytes, pause for %d ms.", p.Name(), p.displayLength, p.Pause)
}
