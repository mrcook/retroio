package types

type BlockType byte

const (
	StandardSpeedData   BlockType = 0x10
	TurboSpeedData      BlockType = 0x11
	PureTone            BlockType = 0x12
	SequenceOfPulses    BlockType = 0x13
	PureData            BlockType = 0x14
	DirectRecording     BlockType = 0x15
	C64RomType          BlockType = 0x16 // deprecated
	C64TurboData        BlockType = 0x17 // deprecated
	CswRecording        BlockType = 0x18
	GeneralizedData     BlockType = 0x19
	PauseTapeCommand    BlockType = 0x20
	GroupStart          BlockType = 0x21
	GroupEnd            BlockType = 0x22
	JumpTo              BlockType = 0x23
	LoopStart           BlockType = 0x24
	LoopEnd             BlockType = 0x25
	CallSequence        BlockType = 0x26
	ReturnFromSequence  BlockType = 0x27
	Select              BlockType = 0x28
	StopTapeWhen48kMode BlockType = 0x2a
	SetSignalLevel      BlockType = 0x2b
	TextDescription     BlockType = 0x30
	Message             BlockType = 0x31
	ArchiveInfo         BlockType = 0x32
	HardwareType        BlockType = 0x33
	EmulationInfo       BlockType = 0x34 // deprecated
	CustomInfo          BlockType = 0x35
	Snapshot            BlockType = 0x40 // deprecated
	GlueBlock           BlockType = 0x5a
)
