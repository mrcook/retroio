package tzx

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tzx/blocks"
	"github.com/mrcook/retroio/spectrum/tzx/blocks/types"
)

// newFromBlockID returns a TZX block based on the type ID byte.
func newFromBlockID(id byte) (Block, error) {
	var block Block

	switch types.BlockType(id) {
	case types.StandardSpeedData:
		block = &blocks.StandardSpeedData{}
	case types.TurboSpeedData:
		block = &blocks.TurboSpeedData{}
	case types.PureTone:
		block = &blocks.PureTone{}
	case types.SequenceOfPulses:
		block = &blocks.SequenceOfPulses{}
	case types.PureData:
		block = &blocks.PureData{}
	case types.DirectRecording:
		block = &blocks.DirectRecording{}
	case types.CswRecording:
		block = &blocks.CswRecording{}
	case types.GeneralizedData:
		block = &blocks.GeneralizedData{}
	case types.PauseTapeCommand:
		block = &blocks.PauseTapeCommand{}
	case types.GroupStart:
		block = &blocks.GroupStart{}
	case types.GroupEnd:
		block = &blocks.GroupEnd{}
	case types.JumpTo:
		block = &blocks.JumpTo{}
	case types.LoopStart:
		block = &blocks.LoopStart{}
	case types.LoopEnd:
		block = &blocks.LoopEnd{}
	case types.CallSequence:
		block = &blocks.CallSequence{}
	case types.ReturnFromSequence:
		block = &blocks.ReturnFromSequence{}
	case types.Select:
		block = &blocks.Select{}
	case types.StopTapeWhen48kMode:
		block = &blocks.StopTapeWhen48kMode{}
	case types.SetSignalLevel:
		block = &blocks.SetSignalLevel{}
	case types.TextDescription:
		block = &blocks.TextDescription{}
	case types.Message:
		block = &blocks.Message{}
	case types.ArchiveInfo:
		block = &blocks.ArchiveInfo{}
	case types.HardwareType:
		block = &blocks.HardwareType{}
	case types.CustomInfo:
		block = &blocks.CustomInfo{}
	case types.GlueBlock:
		// (90 dec, ASCII Letter 'Z')
		block = &blocks.GlueBlock{}
	case types.C64RomType, types.C64TurboData, types.EmulationInfo, types.Snapshot:
		return nil, fmt.Errorf("TZX block ID 0x%02X is deprecated", id)
	default:
		return nil, fmt.Errorf("TZX block ID 0x%02X is not supported", id)
	}
	return block, nil
}
