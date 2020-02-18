// NOTE: unlike the D64/D71, the interleave for a D81 is `1` for both files and directories.
package d81

import "retroio/commodore/disk"

// Directory Track
//
// The directory track should be contained totally on track 40. Sectors 1 and 2
// contain the BAM (Block Availability Map), Sectors 3-39 contain the entries.
//
// Sector 0 holds the disk name and ID.
// The first directory sector is always 40/3, even though the t/s pointer
// (first two bytes) at 40/0 might point somewhere else. It goes linearly up
// the sector count, 3-4-5-6-etc. (an interleave of 1).
// Each sector holds up to eight entries.
type Directory struct {
	HeaderSector Header
	BAM          BlockAvailabilityMap
	Directories  [296]disk.DirectoryFile
}

// Header Sector
//
// The header sector is stored at 40/0, and contains the disk name, ID and DOS
// version bytes, but the BAM is no longer contained here (like the D64).
// NOTE: this is the same size, and a very similar layout to the D64 BAM.
type Header struct {
	//     $00-$01  Track/Sector location of the first directory sector (should
	//             be set to 40/3 but it doesn't matter, and don't trust  what
	//             is there, always go to 40/3 for first directory entry)
	FirstDirTrack  uint8
	FirstDirSector uint8

	// Disk DOS version (see `docs.md` for more information)
	DiskDosVersion byte // Usually 'D' ($44)

	Unused1 uint8 // ($00)

	// 16 character Disk Name (padded with $A0)
	DiskName [16]byte

	Filler1 [2]uint8 // Filled with $A0

	DiskID [2]uint8

	Unknown uint8 // Usually $A0

	// see `docs.md` for more information
	DosVersion  byte // "3"
	DiskVersion byte // "D"

	Filler2 [2]uint8 // Filled with $A0

	Unused2 [227]uint8 // Unused (usually $00)
}

// BAM Layout for the D81
//
// The BAM is different on the 1581 as there are now two separate BAM entries
// (one for each side of the disk), located on Track 40 in Sectors 1 and 2:
// Sector 1: side 1 BAM for tracks  1-40
// Sector 2: side 2 BAM for tracks 41-80
type BlockAvailabilityMap struct {
	// Track/Sector location of next BAM sector.
	//
	// Side 1 BAM should point to $28/$02 (track 40, sector 2)
	// Side 2 BAM has no next sector, and is set to $00/$FF
	NextBamTrack  uint8
	NextBamSector uint8

	// Disk DOS version (see `docs.md` for more information)
	DiskDosVersion uint8 // Usually 'D' ($44)

	// One's complement of version# ($BB)
	DiskDosVersionInverted uint8

	// Disk ID bytes (same as 40/0 Disk ID)
	DiskId [2]uint8

	// I/O byte
	// - bit-7 set   - Verify on
	// - bit-7 clear - Verify off
	// - bit-6 set   - Check header CRC
	// - bit-6 clear - Don't check header CRC
	IO uint8

	// Auto-boot-loader flag (see section at end of document)
	AutoBoot uint8

	// Reserved for future (set to $00)
	Reserved [8]uint8

	// BAM entries for each track on one side of the disk
	// In groups of six bytes, one for the "free sector" count and five for the
	// allocation bitmap.
	BamEntries [40]disk.BamEntry24Bit
}
