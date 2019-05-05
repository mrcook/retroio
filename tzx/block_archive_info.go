package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

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
	TypeID     uint8  // BYTE  Text identification byte
	Length     uint8  // L BYTE  Length of text string
	Characters []byte // CHAR[L] Text string in ASCII format
}

// Headings for the Text ID's.
var headings = map[uint8]string{
	0x00: "Title",     // 00 - Full title
	0x01: "Publisher", // 01 - Software house/publisher
	0x02: "Authors",   // 02 - Author(s)
	0x03: "Year",      // 03 - Year of publication
	0x04: "Language",  // 04 - Language
	0x05: "Category",  // 05 - Game/utility type
	0x06: "Price",     // 06 - Price
	0x07: "Loader",    // 07 - Protection scheme/loader
	0x08: "Origin",    // 08 - Origin
	0xff: "Comment",   // FF - Comment(s)
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (a *ArchiveInfo) Read(file *tape.File) {
	a.Length = file.ReadShort()
	a.StringCount, _ = file.ReadByte()

	for i := 0; i < int(a.StringCount); i++ {
		var t Text
		t.TypeID, _ = file.ReadByte()
		t.Length, _ = file.ReadByte()
		for _, c := range file.ReadBytes(int(t.Length)) {
			t.Characters = append(t.Characters, c)
		}
		a.Strings = append(a.Strings, t)
	}
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (a ArchiveInfo) Id() uint8 {
	return 0x32
}

// Name of the block as given in the TZX specification.
func (a ArchiveInfo) Name() string {
	return "Archive Info"
}

// ToString returns a human readable string of the block data
func (a ArchiveInfo) ToString() string {
	str := ""
	for _, b := range a.Strings {
		// replace newline with comma
		for i := range b.Characters {
			if b.Characters[i] == 0x0a || b.Characters[i] == 0x0d {
				b.Characters[i] = 0x2c // 44d
			}
		}
		str += fmt.Sprintf("  %-10s: %s\n", headings[b.TypeID], b.Characters)
	}

	return str
}
