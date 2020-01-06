package blocks

import (
	"fmt"

	"retroio/spectrum/tap"
	"retroio/spectrum/tzx/blocks/types"
	"retroio/storage"
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
	BlockID     types.BlockType
	Length      uint16 // Length of the whole block (without these two bytes)
	StringCount uint8  // Number of text strings
	Strings     []Text // List of text strings
}

type Text struct {
	TypeID     uint8  // Text identification byte
	Length     uint8  // Length of text string
	Characters []byte // Text string in ASCII format
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
func (a *ArchiveInfo) Read(reader *storage.Reader) error {
	a.BlockID = types.BlockType(reader.ReadByte())
	if a.BlockID != a.Id() {
		return fmt.Errorf("expected block ID 0x%02x, got 0x%02x", a.Id(), a.BlockID)
	}

	a.Length = reader.ReadShort()
	a.StringCount = reader.ReadByte()

	for i := 0; i < int(a.StringCount); i++ {
		var t Text
		t.TypeID = reader.ReadByte()
		t.Length = reader.ReadByte()
		for _, c := range reader.ReadBytes(int(t.Length)) {
			t.Characters = append(t.Characters, c)
		}
		a.Strings = append(a.Strings, t)
	}

	return nil
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (a ArchiveInfo) Id() types.BlockType {
	return types.ArchiveInfo
}

// Name of the block as given in the TZX specification.
func (a ArchiveInfo) Name() string {
	return "Archive Info"
}

func (a ArchiveInfo) BlockData() tap.Block {
	return nil
}

// String returns a human readable string of the block data
// Each character is first converted to a Rune so that Latin characters are preserved.
func (a ArchiveInfo) String() string {
	str := ""
	for _, b := range a.Strings {
		var runes []rune

		for _, c := range b.Characters {
			if c == 0x0a || c == 0x0d {
				// replace newline with comma
				runes = append(runes, rune(0x2c)) // 44d
			} else {
				runes = append(runes, rune(c))
			}
		}

		str += fmt.Sprintf("  %-10s: %s\n", headings[b.TypeID], string(runes))
	}

	return str
}
