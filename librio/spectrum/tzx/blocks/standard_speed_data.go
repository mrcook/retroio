package blocks

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// StandardSpeedData
// ID: 10h (16d)
// This block must be replayed with the standard Spectrum ROM timing values - see the values in
// curly brackets in block ID 11. The pilot tone consists in 8063 pulses if the first data byte
// (flag byte) is < 128, 3223 otherwise. This block can be used for the ROM loading routines AND
// for custom loading routines that use the same timings as ROM ones do.
type StandardSpeedData struct {
	BlockID types.BlockType
	Pause   uint16 // Pause after this block (ms.) {1000}

	// A single .TAP DataBlock consisting of:
	//   WORD    Length of data that follows
	//   BYTE[N] Data as in .TAP files
	DataBlock tap.Block
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (s *StandardSpeedData) Read(reader *storage.Reader) error {
	s.BlockID = types.BlockType(reader.ReadByte())
	if s.BlockID != s.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", s.Id(), s.BlockID)
	}

	s.Pause = reader.ReadShort()

	// Read in the TAP data
	var err error
	tapReader := tap.New(reader)
	length, err := reader.PeekShort()
	if err == nil {
		if length == 19 {
			s.DataBlock, err = tapReader.ReadHeaderBlock()
		} else {
			s.DataBlock, err = tapReader.ReadDataBlock()
		}
	}
	if err != nil {
		return errors.Wrap(err, "unable to read TAP data for StandardSpeedData")
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (s StandardSpeedData) Id() types.BlockType {
	return types.StandardSpeedData
}

// Name of the block as given in the TZX specification.
func (s StandardSpeedData) Name() string {
	return "Standard Speed Data"
}

func (s StandardSpeedData) BlockData() tap.Block {
	return s.DataBlock
}
