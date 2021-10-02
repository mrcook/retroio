package disk

import (
	"bytes"
	"encoding/binary"
)

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
