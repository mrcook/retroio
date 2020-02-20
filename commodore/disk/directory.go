package disk

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// A BAM Entry for each track on the disk
//
// The sector bitmap bytes (three for the D64/D71 and five for D81) indicate
// which sectors are used/free for one track.
//
// A D81 disk with a five byte bitmap gives a total of 40 bits of storage
// (8 bits per byte), one for each sector on the track. If track 40 has these
// five bytes set as $F0, $FF, $2D, $FF, $FE, they would produce a bitmap of:
//
//   F0=11110000, FF=11111111, 2D=00101101, FF=11111111, FE=11111110
//
// The sector bits of *each byte* are stored right-to-left, so to aid in
// understanding the binary notation, flip the bits around.
//
//              111111 11112222 22222233 33333333
//   01234567 89012345 67890123 45678901 23456789
//   -------------------------- -------- --------
//   00001111 11111111 10110100 11111111 01111111
//   ^                     ^             ^
//   sector 0           sector 20     sector 32
//
// The 40 sectors of the first track use bit positions 0-39. If a bit is
// set (1), then the sector is free. Therefore, track 40 has sectors 0-3,
// 17, 20, 22, 23, and 32 in use, the rest are unused.
//
// A D64/D71 would use up to 21 bits from its three bitmap bytes (bits 0-20).
// The bits for non-existent sectors on a track are marked as allocated (reset).
//
// See `docs.md` for lots more information on this set of bytes.
type BamEntry24Bit struct {
	// The number of free sectors on this track
	FreeSectors uint8

	// Bitmap of which sectors are used/free.
	SectorBitmap [3]uint8
}
type BamEntry40Bit struct {
	// The number of free sectors on this track
	FreeSectors uint8

	// Bitmap of which sectors are used/free.
	SectorBitmap [5]uint8
}

// Directory File Entry
//
// All three disk formats (D64, D71, and D81) use the same 32-byte structure,
// with some minor differences with byte $02 (file type), and bytes $1E-$1F
// (file size in sectors).
type DirectoryFile struct {
	// Track/Sector location of next directory sector ($00 $00 if not the first entry in the sector)
	// When the directory is done, the track will be $00, and sector should
	// contain a value of $FF, meaning the whole sector is allocated.
	//
	// Track:
	// - $00 indicates the last sector of the directory
	// - otherwise usually $12 (track 18)
	NextTrack  uint8
	NextSector uint8

	// File type
	// Bit   Description
	// 0-3   The actual file type
	//         Binary Decimal File type
	//          0000     0       DEL
	//          0001     1       SEQ
	//          0010     2       PRG
	//          0011     3       USR
	//          0100     4       REL
	//          0101     5       CBM (partition or sub-directory, only D81)
	//         See `docs.md` for info on illegal values 5-15 (D64/D71) and 6-15 (D81).
	//   4   Unused
	//   5   Used only during SAVE-@ replacement
	//   6   Locked flag (Set produces ">" locked files)
	//   7   Closed flag (Not set produces "*", or "splat" files)
	//
	// Typical values for this location are:
	//   $00 - Scratched (deleted file entry)
	//   $80 - DEL
	//   $81 - SEQ
	//   $82 - PRG
	//   $83 - USR
	//   $84 - REL
	//   $85 - CBM (D81 ONLY)
	FileType uint8

	// Track/sector location of first sector of file
	FirstSectorLocation [2]uint8

	// 16 character filename (in PETASCII, padded with $A0)
	Filename [16]uint8

	// (REL files only)
	// Track/Sector location of first SIDE SECTOR block (D64/D71), or the
	// first SUPER SIDE SECTOR block (D81)
	FirstSideSectorTrack  uint8
	FirstSideSectorSector uint8

	// File record length (REL files only)
	RecordLength uint8 // max. value: 254

	// Unused
	// Except with GEOS disks (the following 6-bytes?)
	Unused [4]uint8

	// Track/sector of replacement file (only used during an @SAVE or an @OPEN command)
	// Note: these may not be used in a D64 - this needs checking.
	ReplacementFileTrack  uint8
	ReplacementFileSector uint8

	// File size in sectors (D64/D71/D81), or partition size in sectors (D81)
	// Little Endian byte order ($1E+$1F*256), with an approx. bytes filesize <= # sectors * 254
	FileSizeInSectors uint16
}

func (f *DirectoryFile) Read(reader *bytes.Reader) error {
	if err := binary.Read(reader, binary.LittleEndian, f); err != nil {
		return err
	}
	return nil
}

func (f DirectoryFile) FileTypeFromID() FileType {
	// Scratch types have a value of 0x00
	if f.FileType == 0x00 {
		return FileType{
			Value:       0xFF,  // not part of the official spec
			Type:        "xxx", // not part of the official spec
			Description: "Scratch File",
		}
	}

	for _, t := range fileTypes {
		if (f.FileType & 0b00000111) == t.Value {
			t.SaveFlag = f.FileType&0b00100000 > 0
			t.LockedFlag = f.FileType&0b01000000 > 0
			t.ClosedFlag = f.FileType&0b10000000 > 0
			return t
		}
	}

	// Should never reach here, but if an unknown type is encountered this will
	// make it more obvious -- not part of the official spec!
	return FileType{
		Value:       0xFF,
		Type:        "???",
		SaveFlag:    false,
		LockedFlag:  false,
		ClosedFlag:  false,
		Description: "Unknown Type",
	}
}

func (f DirectoryFile) PrintableFilename() string {
	name := ""
	for _, c := range f.Filename {
		if c != 0xA0 {
			name += string(c)
		}
	}
	return name
}

type FileType struct {
	Value       uint8
	Type        string
	SaveFlag    bool
	LockedFlag  bool
	ClosedFlag  bool
	Description string
}

func (f FileType) String() string {
	fileType := f.Type

	if !f.ClosedFlag {
		fileType += "*"
	} else if f.LockedFlag {
		fileType += "<"
	}

	// if a scratch file or unknown type: add the description
	if f.Value == 0xFF {
		fileType = fmt.Sprintf("%s (%s)", fileType, f.Description)
	}

	return fileType
}

// Directory File Type Labels
var fileTypes = []FileType{
	{0, "DEL", false, false, true, "Deleted"},
	{1, "SEQ", false, false, true, "Sequential"},
	{2, "PRG", false, false, true, "Program"},
	{3, "USR", false, false, true, "User"},
	{4, "REL", false, false, true, "Relative"},

	// D81 ONLY
	{5, "CBM", false, false, true, "Partition/Sub-directory"},
}

// UsedSectorBitmap is a helper function to determine if a sector is allocated
// in the BAM sector bitmap table. This function will work for both the 24-bit
// and 40-bit maps of the D64, D71, and D81 disk images.
func UsedSectorBitmap(sector uint8, bitmap []uint8) bool {
	// safety check to make sure we don't try accessing a non-existent byte
	if int(sector) > len(bitmap)*8 {
		return false
	}

	if sector <= 7 {
		if reverseBits(bitmap[0])&(1<<sector) > 0 {
			return true
		}
	} else if sector <= 15 {
		sector -= 8
		if reverseBits(bitmap[1])&(1<<sector) > 0 {
			return true
		}
	} else if sector <= 23 {
		sector -= 16
		if reverseBits(bitmap[2])&(1<<sector) > 0 {
			return true
		}
	} else if sector <= 31 {
		sector -= 24
		if reverseBits(bitmap[3])&(1<<sector) > 0 {
			return true
		}
	} else if sector <= 39 {
		sector -= 32
		if reverseBits(bitmap[4])&(1<<sector) > 0 {
			return true
		}
	}

	return false
}

func reverseBits(b byte) byte {
	var d byte
	for i := 0; i < 8; i++ {
		d <<= 1
		d |= b & 1
		b >>= 1
	}
	return d
}
