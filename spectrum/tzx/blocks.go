package tzx

import (
	"fmt"
	blocks2 "retroio/spectrum/tzx/blocks"
	types2 "retroio/spectrum/tzx/blocks/types"

	"retroio/tape"
)

// newFromBlockID returns a TZX block based on the type ID byte.
func newFromBlockID(id byte) (tape.Block, error) {
	var block tape.Block

	switch types2.BlockType(id) {
	case types2.StandardSpeedData:
		block = &blocks2.StandardSpeedData{}
	case types2.TurboSpeedData:
		block = &blocks2.TurboSpeedData{}
	case types2.PureTone:
		block = &blocks2.PureTone{}
	case types2.SequenceOfPulses:
		block = &blocks2.SequenceOfPulses{}
	case types2.PureData:
		block = &blocks2.PureData{}
	case types2.DirectRecording:
		block = &blocks2.DirectRecording{}
	case types2.CswRecording:
		block = &blocks2.CswRecording{}
	case types2.GeneralizedData:
		block = &blocks2.GeneralizedData{}
	case types2.PauseTapeCommand:
		block = &blocks2.PauseTapeCommand{}
	case types2.GroupStart:
		block = &blocks2.GroupStart{}
	case types2.GroupEnd:
		block = &blocks2.GroupEnd{}
	case types2.JumpTo:
		block = &blocks2.JumpTo{}
	case types2.LoopStart:
		block = &blocks2.LoopStart{}
	case types2.LoopEnd:
		block = &blocks2.LoopEnd{}
	case types2.CallSequence:
		block = &blocks2.CallSequence{}
	case types2.ReturnFromSequence:
		block = &blocks2.ReturnFromSequence{}
	case types2.Select:
		block = &blocks2.Select{}
	case types2.StopTapeWhen48kMode:
		block = &blocks2.StopTapeWhen48kMode{}
	case types2.SetSignalLevel:
		block = &blocks2.SetSignalLevel{}
	case types2.TextDescription:
		block = &blocks2.TextDescription{}
	case types2.Message:
		block = &blocks2.Message{}
	case types2.ArchiveInfo:
		block = &blocks2.ArchiveInfo{}
	case types2.HardwareType:
		block = &blocks2.HardwareType{}
	case types2.CustomInfo:
		block = &blocks2.CustomInfo{}
	case types2.GlueBlock:
		// (90 dec, ASCII Letter 'Z')
		block = &blocks2.GlueBlock{}
	case types2.C64RomType, types2.C64TurboData, types2.EmulationInfo, types2.Snapshot:
		return nil, fmt.Errorf("TZX block ID 0x%02X is deprecated", id)
	default:
		return nil, fmt.Errorf("TZX block ID 0x%02X is not supported", id)
	}
	return block, nil
}
