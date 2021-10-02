package blocks

import (
	"log"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	gdata "github.com/mrcook/retroio/librio/spectrum/tzx/blocks/generalized_data"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// GeneralizedData
// ID: 19h (25d)
// This block has been specifically developed to represent an extremely wide range of data encoding techniques.
// The basic idea is that each loading component (pilot tone, sync pulses, data) is associated to a
// specific sequence of pulses, where each sequence (wave) can contain a different number of pulses
// from the others. In this way we can have a situation where bit 0 is represented with 4 pulses and
// bit 1 with 8 pulses.
type GeneralizedData struct {
	BlockID      types.BlockType
	Length       uint32           // Block length (without these four bytes)
	Pause        uint16           // Pause after this block (ms)
	TOTP         uint32           // Total number of symbols in pilot/sync block (can be 0)
	NPP          uint8            // Maximum number of pulses per pilot/sync symbol
	ASP          uint8            // Number of pilot/sync symbols in the alphabet table (0=256)
	TOTD         uint32           // Total number of symbols in data stream (can be 0)
	NPD          uint8            // Maximum number of pulses per data symbol
	ASD          uint8            // Number of data symbols in the alphabet table (0=256)
	PilotSymbols []gdata.Symbol   // 0x12  SYMDEF[ASP] Pilot and sync symbols definition table: this field is present only if TOTP>0
	PilotStreams []gdata.PilotRLE // 0x12+ (2*NPP+1)*ASP - PRLE[TOTP]  Pilot and sync data stream: this field is present only if TOTP>0
	DataSymbols  []gdata.Symbol   // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+TOTP*3  - SYMDEF[ASD] Data symbols definition table: this field is present only if TOTD>0
	DataStreams  []uint8          // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+ TOTP*3+(2*NPD+1)*ASD - BYTE[DS]  Data stream: this field is present only if TOTD>0
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (g *GeneralizedData) Read(reader *storage.Reader) error {
	log.Fatal("GeneralizedData is not currently supported - unable to continue.")

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (g GeneralizedData) Id() types.BlockType {
	return types.GeneralizedData
}

// Name of the block as given in the TZX specification.
func (g GeneralizedData) Name() string {
	return "Generalized Data"
}

func (g GeneralizedData) BlockData() tap.Block {
	return nil
}
