package d81

import (
	"bytes"
	"encoding/binary"

	"github.com/mrcook/retroio/librio/commodore/disk"
)

// Header Sector
//
// The header sector is stored at 40/0, and contains the disk name, ID and DOS
// version bytes, but the BAM is no longer contained here (like the D64).
// NOTE: this is the same size, and a very similar layout to the D64 BAM.
type Header struct {
	// Track/Sector location of the first directory sector (should
	// be set to 40/3 but it doesn't matter, and don't trust what
	// is there, always go to 40/3 for first directory entry)
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

func (h *Header) Read(sector disk.Sector) error {
	reader := bytes.NewReader(sector.Data[:])
	if err := binary.Read(reader, binary.LittleEndian, h); err != nil {
		return err
	}

	return nil
}
