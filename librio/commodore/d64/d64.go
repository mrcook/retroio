// Package d64 implements reading Commodore emulator D64 disk images.
//
// A D64 is [basically] a sector-for-sector copy of a 1540/1541 disk.
// The original media (a 5.25" disk) had a variable number of 256-byte sectors
// per track and data densities across the disk to optimize available storage.
//
// Track counting starts at 1 (not 0) and goes up to 35. Sector counting starts at 0.
// It's only possible to know the D64 type being read by first checking its file size.
//
// Additional D64 documentation can be found in the `docs.md` file.
package d64

import (
	"fmt"

	"github.com/mrcook/retroio/librio/commodore/disk"
	"github.com/mrcook/retroio/librio/storage"
)

const (
	DirectoryTrackNumber = 18 - 1 // tracks start counting from 1 not zero
	DirEntriesPerSector  = 8
	BamSectorBitmapSize  = 3 // 24 bit BAM bitmap
)

type D64 struct {
	disk *disk.Disk
	cbm  *Directory
}

var Geometry = []disk.TrackGeometry{
	{0, 1, 17, 21, 357},
	{0, 18, 24, 19, 133},
	{0, 25, 30, 18, 108},
	{0, 31, 35, 17, 85},
	{0, 36, 40, 17, 85},
	{0, 41, 42, 17, 34},
}

func (d *D64) Read(reader *storage.Reader) error {
	if err := d.disk.Read(reader); err != nil {
		return fmt.Errorf("error reading the disk: %w", err)
	}

	d.cbm = &Directory{}
	if err := d.cbm.Read(reader); err != nil {
		return fmt.Errorf("error reading the directory: %w", err)
	}

	return nil
}
