// Package contains definitions for all the TZX data blocks
package tape

type Block interface {
	Id() int
	Name() string
	Metadata() string
}

// Block ID table, as `HEX (DECIMAL)` values
//
// ID 10 (16): StandardSpeedData
// ID 11 (17): TurboSpeedData
// ID 12 (18): PureTone
// ID 13 (19): SequenceOfPulses
// ID 14 (20): PureData
// ID 15 (21): DirectRecording
// ID 18 (24): CswRecording
// ID 19 (25): GeneralizedData
// ID 20 (32): PauseTheTapeCommand
// ID 21 (33): GroupStart
// ID 22 (34): GroupEnd
// ID 23 (35): JumpTo
// ID 24 (36): LoopStart
// ID 25 (37): LoopEnd
// ID 26 (38): CallSequence
// ID 27 (39): ReturnFromSequence
// ID 28 (40): Select
// ID 2A (42): StopTheTapeIfIn48kMode
// ID 2B (43): SetSignalLevel
// ID 30 (48): TextDescription
// ID 31 (49): Message
// ID 32 (50): ArchiveInfo
// ID 33 (51): HardwareType
// ID 35 (53): CustomInfo
// ID 5A (90): GlueBlock
