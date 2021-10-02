package d64

import (
	"bytes"
	"encoding/binary"

	"github.com/mrcook/retroio/librio/commodore/disk"
)

// BAM Layout for the D64 and D71 disks.
//
// D64: BAM information is stored on track 18/0
// D71: BAM is different on the 1571 as it must now manage the 35 extra tracks.
//      Most of the extra BAM information is stored in sector 0 on track 53
//      (the equivalent of track 18, but on the flip side of the disk), with
//      the remaining sectors on track 53 marked as allocated.
//      This makes the remainder of track 53 unusable.
type BlockAvailabilityMap struct {
	// Track/Sector location of the first directory sector - should be set to
	// 18/1 but it doesn't matter, and don't trust what is there, always go
	// to 18/1 for first directory entry.
	FirstDirTrack  uint8
	FirstDirSector uint8

	// Disk DOS version (see `docs.md` for more information)
	DiskDosVersion byte // Usually 'A' ($41)

	// D64: Unused ($00)
	// D71: Double-sided flag
	//      $00 - Single sided disk
	//      $80 - Double sided disk
	NumberOfSidesFlag uint8

	// BAM entries for each track
	// In groups of four bytes per track, starting on track 1.
	Entries [35]disk.BamEntry24Bit

	// Disk Name (padded with $A0)
	DiskName [16]byte

	Filler1 [2]uint8 // Filled with $A0

	DiskID [2]byte

	Unknown uint8 // Usually $A0

	// see `docs.md` for more information
	DosVersion  byte // "2"
	DiskVersion byte // "A"

	Filler2 [4]uint8 // Filled with $A0

	// The following 85 bytes are unused on the standard D64 disks,
	// and are filled with $00 values.
	Unused2 uint8

	// DOLPHIN DOS: track 36-40 BAM entries
	DolphinDosBAM [20]uint8

	// SPEED DOS: track 36-40 BAM entries
	SpeedDosBAM [20]uint8

	Unused3 [9]uint8

	// D71 ONLY
	// Free sector count for the 35 tracks on side 2 of a 1571 disk.
	// These 1-byte per track "free sector" entries for tracks 36-70 are
	// likely included here due to memory restrictions in the 1571 drive.
	// NOTE: in a D64 these values would normally be with the 4-byte BAM entry,
	// but the rest of the entry is contained on 53/0.
	FreeSectorCount [35]uint8
}

func (b *BlockAvailabilityMap) Read(sector disk.Sector) error {
	reader := bytes.NewReader(sector.Data[:])
	if err := binary.Read(reader, binary.LittleEndian, b); err != nil {
		return err
	}

	return nil
}
