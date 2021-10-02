package blocks

import (
	"log"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// Snapshot (deprecated in v1.20)
// ID: 40h (64d)
//
// This would enable one to snapshot the game at the start and still have all the
// tape blocks (level data, etc.) in the same file. Only `.Z80` and `.SNA`
// snapshots are supported for compatibility reasons! The emulator should take
// care of that the snapshot is not taken while the actual Tape loading is taking
// place (which doesn't do much sense). And when an emulator encounters the
// snapshot block it should load it and then continue with the next block.
type Snapshot struct {
	blockID types.BlockType

	// Snapshot format type:
	//   0x00: .Z80 format
	//   0x01: .SNA format
	formatType uint8

	length [3]uint8 // Snapshot length (24-bit)
	data   []uint8  // Snapshot data itself
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (d *Snapshot) Read(reader *storage.Reader) error {
	log.Fatal("Snapshot is not currently supported - unable to continue.")

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (d Snapshot) Id() types.BlockType {
	return types.Snapshot
}

// Name of the block as given in the TZX specification.
func (d Snapshot) Name() string {
	return "Snapshot"
}

func (d Snapshot) BlockData() tap.Block {
	return nil
}
