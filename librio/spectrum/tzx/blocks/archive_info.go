package blocks

import (
	"fmt"

	"github.com/mrcook/retroio/librio/spectrum/tap"
	"github.com/mrcook/retroio/librio/spectrum/tzx/blocks/types"
	"github.com/mrcook/retroio/librio/storage"
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
	Length      uint16      // Length of the whole block (without these two bytes)
	StringCount uint8       // Number of text strings
	Strings     []TextField // List of text strings
}

type TextField struct {
	TypeID     FieldType // Text identification byte
	Length     uint8     // Length of text string
	Characters []byte    // Text string in ASCII format TODO: use string type?
}

type FieldType uint8

// Types IDs for the archive info fields.
const (
	Title     FieldType = 0x00 // Full title
	Publisher FieldType = 0x01 // Software house/publisher
	Authors   FieldType = 0x02 // Author(s)
	Year      FieldType = 0x03 // Year of publication
	Language  FieldType = 0x04 // Language
	Category  FieldType = 0x05 // Game/utility type (e.g. Arcade, Puzzle)
	Price     FieldType = 0x06 // Price (including currency: GBP, DEM, FRF, EUR, etc.)
	Loader    FieldType = 0x07 // Protection scheme/loader (e.g. SpeedLock 1)
	Origin    FieldType = 0x08 // Origin (e.g. Original, Budget re-release)
	Comment   FieldType = 0xFF // Comment(s)
)

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
		var t TextField
		t.TypeID = FieldType(reader.ReadByte())
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
