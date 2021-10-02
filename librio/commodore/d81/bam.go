package d81

import (
	"bytes"
	"encoding/binary"

	"github.com/mrcook/retroio/librio/commodore/disk"
)

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
	// Entries [40]disk.BamEntry40Bit // TODO: fix me
}

func (b *BlockAvailabilityMap) Read(sector disk.Sector) error {
	reader := bytes.NewReader(sector.Data[:])
	if err := binary.Read(reader, binary.LittleEndian, b); err != nil {
		return err
	}

	return nil
}
