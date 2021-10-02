package trd

// Specification (256 bytes) stores the disk information obtained from the
// 9th sector of the system track.
type Specification struct {
	EndOfDirectory uint8      // must be 0 to indicate the end of the directory
	Unused1        [224]uint8 // Unused, filled with 0x00

	// First free sector (0..15)
	// Representing the logical sector number of the next free sector on the disk.
	// If the disk is full, the sector number is = count of sectors.
	NextFreeSector uint8

	// First free track: logical sector=first free track*16+first free sector
	NextFreeTrack uint8

	// Disk Type (0x16..0x19)
	//   22: double-sided, 80 tracks
	//   23: double-sided, 40 tracks
	//   24: single-sided, 80 tracks
	//   25: single-sided, 40 tracks
	DiskType DiskType

	NumFiles       uint8  // 0..128; the count of non-deleted files
	NumFreeSectors uint32 // Number of free sectors on disk
	TrdId          uint8  // Always 0x10 - TR-DOS ID
	Unused2        uint8  // Unused: filled with 0x00

	// Unused and filled with spaces (0x20),
	// or a disk protecting password filled up with spaces (thanks to Németh Zoltán Gábor for this info)
	Unused3 [9]uint8

	Unused4         uint8    // Unused: filled with 0x00
	NumDeletedFiles uint8    // Number of deleted files on disk
	Label           [8]uint8 // Disk label
	Unused5         [3]uint8 // Unused, filled with 0x00
}
