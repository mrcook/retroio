package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
)

// TODO: I'd like to use generics instead of this interface
// Block is an interface for Tape data blocks
type Block interface {
	Read(reader *storage.Reader) error
	Id() types.BlockType
	Name() string
	BlockData() tap.Block
}

// NewFromBlockID returns a TZX block based on the type ID byte.
func NewFromBlockID(id byte) (Block, error) {
	var block Block

	switch types.BlockType(id) {
	case types.StandardSpeedData:
		block = &StandardSpeedData{}
	case types.TurboSpeedData:
		block = &TurboSpeedData{}
	case types.PureTone:
		block = &PureTone{}
	case types.SequenceOfPulses:
		block = &SequenceOfPulses{}
	case types.PureData:
		block = &PureData{}
	case types.DirectRecording:
		block = &DirectRecording{}
	case types.CswRecording:
		block = &CswRecording{}
	case types.GeneralizedData:
		block = &GeneralizedData{}
	case types.PauseTapeCommand:
		block = &PauseTapeCommand{}
	case types.GroupStart:
		block = &GroupStart{}
	case types.GroupEnd:
		block = &GroupEnd{}
	case types.JumpTo:
		block = &JumpTo{}
	case types.LoopStart:
		block = &LoopStart{}
	case types.LoopEnd:
		block = &LoopEnd{}
	case types.CallSequence:
		block = &CallSequence{}
	case types.ReturnFromSequence:
		block = &ReturnFromSequence{}
	case types.Select:
		block = &Select{}
	case types.StopTapeWhen48kMode:
		block = &StopTapeWhen48kMode{}
	case types.SetSignalLevel:
		block = &SetSignalLevel{}
	case types.TextDescription:
		block = &TextDescription{}
	case types.Message:
		block = &Message{}
	case types.ArchiveInfo:
		block = &ArchiveInfo{}
	case types.HardwareType:
		block = &HardwareType{}
	case types.CustomInfo:
		block = &CustomInfo{}
	case types.GlueBlock:
		// (90 dec, ASCII Letter 'Z')
		block = &GlueBlock{}
	case types.C64RomType, types.C64TurboData, types.EmulationInfo, types.Snapshot:
		return nil, fmt.Errorf("TZX block ID 0x%02X is deprecated", id)
	default:
		return nil, fmt.Errorf("TZX block ID 0x%02X is not supported", id)
	}
	return block, nil
}
