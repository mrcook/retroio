package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/hardware"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// HardwareType
// ID: 33h (51d)
// This blocks contains information about the hardware that the programs on this tape use.
// Please include only machines and hardware for which you are 100% sure that it either runs
// (or doesn't run) on or with, or you know it uses (or doesn't use) the hardware or special
// features of that machine.
// If the tape runs only on the ZX81 (and TS1000, etc.) then it clearly won't work on any
// Spectrum or Spectrum variant, so there's no need to list this information.
// If you are not sure or you haven't tested a tape on some particular machine/hardware
// combination then do not include it in the list.
type HardwareType struct {
	BlockID   types.BlockType
	TypeCount uint8           // Number of machines and hardware types for which info is supplied
	Machines  []hardware.Info // List of machines and hardware
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (h *HardwareType) Read(reader *storage.Reader) error {
	h.BlockID = types.BlockType(reader.ReadByte())
	if h.BlockID != h.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", h.Id(), h.BlockID)
	}

	h.TypeCount = reader.ReadByte()

	for i := 0; i < int(h.TypeCount); i++ {
		var m hardware.Info
		m.Type = hardware.HwType(reader.ReadByte())
		m.Id = reader.ReadByte()
		m.Compatibility = hardware.Compatibility(reader.ReadByte())
		h.Machines = append(h.Machines, m)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (h HardwareType) Id() types.BlockType {
	return types.HardwareType
}

// Name of the block as given in the TZX specification.
func (h HardwareType) Name() string {
	return "Hardware"
}

func (h HardwareType) BlockData() tap.Block {
	return nil
}
