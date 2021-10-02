package blocks

import (
	"log"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// Emulation info (deprecated in v1.20)
// ID: 34h (52d)
//
// This is a special block that would normally be generated only by emulators.
// For now it contains info on everything I could find that other formats
// support. Please inform me of any additions/corrections since this is a very
// important part for emulators.
//
// Those bits that are not used by the emulator that stored the info, should be
// left at their DEFAULT values.
type EmulationInfo struct {
	BlockID types.BlockType

	// General emulation flags.
	//   bit 0   : R-register emulation [1] bit 1 : `LDIR` emulation [1]
	//   bit 2   : high resolution colour emulation with true interrupt freq. [1]
	//   bit 3,4 : video synchronisation: 1=high, 3=low, 0,2=normal [0]
	//   bit 5   : fast loading when ROM load routine is used [1]
	//   bit 6   : border emulation [1]
	//   bit 7   : screen refresh mode (1: ON, 0: OFF) [1]
	//   bit 8   : start playing the tape immediately [0]
	//             If this is 0 then the emulator should only load the info
	//             blocks and WAIT when it encounters first DATA block
	//   bit 9   : auto type `LOAD""` or press ENTER when in 128k mode [0]
	Flags uint16

	// Screen refresh delay.
	//   1 - 255 (interrupts between refreshes) [1]
	//   (used when screen refresh mode is ON)
	RefreshDelay uint8

	InterruptFrequency uint16   // Interrupt Frequency : 0 - 999 Hz
	Reserved           [3]uint8 // Reserved for future expansion
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (d *EmulationInfo) Read(reader *storage.Reader) error {
	log.Fatal("EmulationInfo is not currently supported - unable to continue.")

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (d EmulationInfo) Id() types.BlockType {
	return types.EmulationInfo
}

// Name of the block as given in the TZX specification.
func (d EmulationInfo) Name() string {
	return "Emulation Info"
}

func (d EmulationInfo) BlockData() tap.Block {
	return nil
}
