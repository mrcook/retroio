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

// GroupStart
// ID: 21h (33d)
// This block marks the start of a group of blocks which are to be treated as one single
// (composite) block. This is very handy for tapes that use lots of subblocks like Bleepload
// (which may well have over 160 custom loading blocks). You can also give the group a name
// (example 'Bleepload Block 1').
// For each group start block, there must be a group end block. Nesting of groups is not allowed.
type GroupStart struct {
	Length uint8  // L BYTE  Length of the group name string
	Name   []byte // CHAR[L] Group name in ASCII format (please keep it under 30 characters long)
}

// GroupEnd
// ID: 22h (34d)
// This indicates the end of a group. This block has no body.
type GroupEnd struct{}

// JumpTo
// ID: 23h (35d)
// This block will enable you to jump from one block to another within the file. The value is a
// signed short word (usually 'signed short' in C); Some examples:
//   - Jump 0 = 'Loop Forever' - this should never happen
//   - Jump 1 = 'Go to the next block' - it is like NOP in assembler ;)
//   - Jump 2 = 'Skip one block'
//   - Jump -1 = 'Go to the previous block'
// All blocks are included in the block count!
type JumpTo struct {
	Value int16 // WORD  Relative jump value
}

// LoopStart
// ID: 24h (36d)
// If you have a sequence of identical blocks, or of identical groups of blocks, you can use
// this block to tell how many times they should be repeated. This block is the same as the
// FOR statement in BASIC.
type LoopStart struct {
	RepetitionCount uint16 // WORD  Number of repetitions (greater than 1)
}

// LoopEnd
// ID: 25h (37d)
// This is the same as BASIC's NEXT statement. It means that the utility should jump back to the
// start of the loop if it hasn't been run for the specified number of times.
// This block has no body.
type LoopEnd struct{}

// CallSequence
// ID: 26h (38d)
// This block is an analogue of the CALL Subroutine statement. It basically executes a sequence of
// blocks that are somewhere else and then goes back to the next block. Because more than one call
// can be normally used you can include a list of sequences to be called. The 'nesting' of call
// blocks is also not allowed for the simplicity reasons. You can, of course, use the CALL blocks
// in the LOOP sequences and vice versa. The value is relative for the obvious reasons - so that
// you can add some blocks in the beginning of the file without disturbing the call values. Please
// take a look at 'Jump To Block' for reference on the values.
type CallSequence struct {
	Count  uint16   // N WORD  Number of calls to be made
	Blocks []uint16 // WORD[N] Array of call block numbers (relative-signed offsets)
}

// ReturnFromSequence
// ID: 27h (39d)
// This block indicates the end of the Called Sequence. The next block played will be the block after
// the last CALL block (or the next Call, if the Call block had multiple calls).
// This block has no body.
type ReturnFromSequence struct{}

// Select
// ID: 28h (40d)
// This block is useful when the tape consists of two or more separately-loadable parts. With this
// block, you are able to select one of the parts and the utility/emulator will start loading from
// that block. For example you can use it when the game has a separate Trainer or when it is a
// multi-load. Of course, to make some use of it the emulator/utility has to show a menu with the
// selections when it encounters such a block. All offsets are relative signed words.
type Select struct {
	Length     uint16      // WORD  Length of the whole block (without these two bytes)
	Count      uint8       // N BYTE  Number of selections
	Selections []Selection // SELECT[N] List of selections
}

type Selection struct {
	RelativeOffset uint16 // WORD  Relative Offset
	Length         uint8  // L BYTE  Length of description text
	Description    uint8  // CHAR[L] Description text (please use single line and max. 30 chars)
}

// StopTheTapeIfIn48kMode
// ID: 2Ah (42d)
// When this block is encountered, the tape will stop ONLY if the machine is an 48K Spectrum.
// This block is to be used for multi-loading games that load one level at a time in 48K mode,
// but load the entire tape at once if in 128K mode.
// This block has no body of its own, but follows the extension rule.
type StopTheTapeIfIn48kMode struct {
	Length uint32 // DWORD Length of the block without these four bytes (0)
}

// SetSignalLevel
// ID: 2Bh (43d)
// This block sets the current signal level to the specified value (high or low). It should be used
// whenever it is necessary to avoid any ambiguities, e.g. with custom loaders which are level-sensitive.
type SetSignalLevel struct {
	Length      uint32 // DWORD Block length (without these four bytes)
	SignalLevel uint8  // BYTE  Signal level (0=low, 1=high)
}

// TextDescription
// ID: 30h (48d)
// This is meant to identify parts of the tape, so you know where level 1 starts, where to rewind
// to when the game ends, etc. This description is not guaranteed to be shown while the tape is
// playing, but can be read while browsing the tape or changing the tape pointer.
// The description can be up to 255 characters long but please keep it down to about 30 so the
// programs can show it in one line (where this is appropriate).
type TextDescription struct {
	Length      uint8  // N BYTE  Length of the text description
	Description []byte // CHAR[N] Text description in ASCII format
}

// Message
// ID: 31h (49d)
// This will enable the emulators to display a message for a given time. This should not stop the
// tape and it should not make silence. If the time is 0 then the emulator should wait for the user
// to press a key.
// The text message should:
//   - stick to a maximum of 30 chars per line;
//   - use single 0x0D (13 decimal) to separate lines;
//   - stick to a maximum of 8 lines.
type Message struct {
	DisplayTime uint8  // BYTE  Time (in seconds) for which the message should be displayed
	Length      uint8  // N BYTE  Length of the text message
	Message     []byte // CHAR[N] Message that should be displayed in ASCII format
}

// ArchiveInfo
// ID: 32h (50d)
// Use this block at the beginning of the tape to identify the title of the game, author,
// publisher, year of publication, price (including the currency), type of software (arcade
// adventure, puzzle, word processor, ...), protection scheme it uses (Speedlock 1, Alkatraz,
// ...) and its origin (Original, Budget re-release, ...), etc. This block is built in a way
// that allows easy future expansion. The block consists of a series of text strings. Each text
// has its identification number (which tells us what the text means) and then the ASCII text.
// To make it possible to skip this block, if needed, the length of the whole block is at the
// beginning of it.
// If all texts on the tape are in English language then you don't have to supply the 'Language' field.
type ArchiveInfo struct {
	Length      uint16 // WORD  Length of the whole block (without these two bytes)
	StringCount uint8  // N BYTE  Number of text strings
	Strings     []Text // TEXT[N] List of text strings
}

type Text struct {
	TypeID uint8 // BYTE  Text identification byte:
	//                       00 - Full title
	//                       01 - Software house/publisher
	//                       02 - Author(s)
	//                       03 - Year of publication
	//                       04 - Language
	//                       05 - Game/utility type
	//                       06 - Price
	//                       07 - Protection scheme/loader
	//                       08 - Origin
	//                       FF - Comment(s)
	Length     uint8  // L BYTE  Length of text string
	Characters []byte // CHAR[L] Text string in ASCII format
}

// HardwareType
// ID: 33h (51d)
// This blocks contains information about the hardware that the programs on this tape use.
// Please include only machines and hardware for which you are 100% sure that it either runs
// (or doesn't run) on or with, or you know it uses (or doesn't use) the hardware or special
// features of that machine.
type HardwareType struct {
	TypeCount uint8          // N BYTE     Number of machines and hardware types for which info is supplied
	Machines  []HardwareInfo // HWINFO[N]  List of machines and hardware
}

// HardwareInfo
// A list of hardware types and ID and sub-types can be found in the TZX specification
// (https://www.worldofspectrum.org/TZXformat.html)
type HardwareInfo struct {
	Type        uint8 // BYTE  Hardware type
	Id          uint8 // BYTE  Hardware ID
	Information uint8 // BYTE  Hardware information:
	//                           00 - The tape RUNS on this machine or with this hardware,
	//                                but may or may not use the hardware or special features of the machine.
	//                           01 - The tape USES the hardware or special features of the machine,
	//                                such as extra memory or a sound chip.
	//                           02 - The tape RUNS but it DOESN'T use the hardware
	//                                or special features of the machine.
	//                           03 - The tape DOESN'T RUN on this machine or with this hardware.
}

// CustomInfo
// ID: 35h (53d)
// This block can be used to save any information you want. For example, it might contain some
// information written by a utility, extra settings required by a particular emulator, or even
// poke data.
type CustomInfo struct {
	Identification []byte  // CHAR[10]  Identification string (in ASCII)
	Length         uint32  // L DWORD   Length of the custom info
	Info           []uint8 // BYTE[L]   Custom info
}

// GlueBlock
// ID: 5Ah (90d)
// This block is generated when you merge two ZX Tape files together. It is here so that you can
// easily copy the files together and use them. Of course, this means that resulting file would
// be 10 bytes longer than if this block was not used. All you have to do if you encounter this
// block ID is to skip next 9 bytes.
// If you can avoid using this block for this purpose, then do so; it is preferable to use a
// utility to join the two files and ensure that they are both of the higher version number.
type GlueBlock struct {
	Value [9]byte // BYTE[9] Value: { "XTape!",0x1A,MajR,MinR } Just skip these 9 bytes and you will end up on the next ID.
}
