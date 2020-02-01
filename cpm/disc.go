// Package cpm implements reading CP/M compatible disc image files.
package cpm

import (
	"fmt"

	"retroio/storage"
)

const (
	RecordSize = 128 // CP/M records are 128 bytes in length
)

type Disc struct {
	reader *storage.Reader

	DPH DiskParameterHeader
	DPB DiskParameterBlock

	// Which of these live on the DPH/DPB?
	//
	// Drive Characteristics
	ReservedTracks     uint8
	SectorsPerTrack    uint8
	RecordsPerBlock    uint8
	RecordsPerExtent   uint8
	CheckedDirectories uint8
	DirectoryEntries   uint8  // 32 Byte Directory Entries
	DriveCapacity      uint16 // Kilobyte Drive Capacity
	RecordCapacity     uint16 // 128 Byte Record Capacity
	// Characteristics sizes
	SectorSize           uint16 // Sector size in bytes
	TrackCount           uint16 // Number of tracks
	SectorCount          uint16 // Number of sectors
	BlockSize            uint16 // Block size
	DirectoryCount       uint16 // Number of directory entries
	LogicalSectorSkew    uint16 // Logical sector skew
	ReservedSystemTracks uint8  // Number of reserved system tracks (optional)
	Offset               uint16 // Offset to start of volume (optional)

	Tracks []Track
}

type DiskParameterHeader interface{}
type DiskParameterBlock interface{}

type Track struct {
	Number  uint8
	Sectors []Sector
}

type Sector struct {
	Number     uint8
	Size       uint16 // Needed?
	PageNumber uint8  // Needed?
	ModuleName string // Needed?
}

func New(reader *storage.Reader) *Disc {
	return &Disc{reader: reader}
}

func (d *Disc) Read() error {
	return nil
}

// DisplayGeometry prints the disk, track and sector metadata to the terminal.
func (d Disc) DisplayGeometry() {
	fmt.Println("DISC INFORMATION:")
}

// CP/M File Control Block
// Reference: http://www.seasip.info/Cpm/fcb.html
//
// The File Control Block is a 36-byte data structure (33 bytes in CP/M 1).
//
// Emulators at BDOS level need to be aware of how CP/M uses the bytes EX, S2,
// and CR. Some programs manipulate these bytes to perform "seek" operations in
// files without using the random-access calls.
//
// CR = current record,   ie (file pointer % 16384)  / 128
// EX = current extent,   ie (file pointer % 524288) / 16384
// S2 = extent high byte, ie (file pointer / 524288).
// The CP/M Plus source code refers to this use of the S2 byte as 'module number'.
//
// DR F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 EX S1 S2 RC  .FILENAMETYP...
// AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL  ...............
// CR R0 R1 R2                                      ....
type FileControlBlock struct {
	// DR - Drive.
	//     0: default
	//  1-16: for A-P
	//
	// In DOSPLUS, bit 7 can be set to indicate that the operation should work
	// with sub-directories rather than files.
	Drive uint8

	// Fn - Filename, 7-bit ASCII.
	// The top bits of the filename bytes (usually referred to as F1' to F8')
	// have the following meanings:
	// F1'-F4' - User-defined attributes.
	//           Any program can use them in any way it likes. The filename in
	//           the disc directory has the corresponding bits set.
	// F5'-F8' - Interface attributes.
	//           They modify the behaviour of various BDOS functions or
	//           indicate error conditions. In the directory these bits are
	//           always zero.
	Filename [8]uint8

	// Tn - Filetype, 7-bit ASCII.
	// T1' - Read-Only.
	// T2' - System (hidden). System files in user 0 can be opened from other user areas.
	// T3' - Archive. Set if the file has not been changed since it was last copied.
	FileType [3]uint8

	// EX - Set this to 0 when opening a file and then leave it to CP/M.
	// You can rewind a file by setting EX, RC, S2 and CR to 0.
	Extent uint8

	// S1 - Reserved.
	S1 uint8

	// S2 - Reserved.
	S2 uint8

	// RC - Set this to 0 when opening a file and then leave it to CP/M.
	RecordCount uint8

	// AL - Allocation.
	// Image of the second half of the directory entry, containing the file's
	// allocation (which disc blocks it owns).
	Allocation [16]uint8

	// CR - Current record within extent.
	// It is usually best to set this to 0 immediately after a file has been
	// opened and then ignore it.
	CurrentRecord uint8

	// Rn - Random access record number (not CP/M 1).
	// A 16-bit value in CP/M 2 (with R2 used for overflow); an 18-bit value in CP/M 3.
	RndRecordNumber [3]uint8
}
