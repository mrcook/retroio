// Package d71 implements reading Commodore D71 image files.
//
// Similar to the D64 (1541), the 1571 drive can operate in either single-sided
// (1541 compatible) mode or double-sided (1571) mode.
//
// The D71 is a sector-for-sector copy of a 1571 disk, with 70 tracks, double
// that of the 1541, making a DOS file size of 349696 bytes. If the error byte
// block (1366 bytes) is attached, the file size wil be 351062 bytes.
//
// Additional D71 documentation can be found in the `docs.md` file.
package d71

import (
	"fmt"

	"github.com/mrcook/retroio/librio/commodore/disk"
	"github.com/mrcook/retroio/librio/storage"
)

const (
	DirectoryTrackNumber         = 18 - 1 // tracks start counting from 1 not zero
	DirEntriesPerSector          = 8
	BamSectorBitmapSize          = 3      // 24 bit BAM bitmap
	DirectoryBamExtraTrackNumber = 53 - 1 // tracks start counting from 1 not zero
)

var Geometry = []disk.TrackGeometry{
	{0, 1, 17, 21, 357},
	{0, 18, 24, 19, 133},
	{0, 25, 30, 18, 108},
	{0, 31, 35, 17, 85},
	{1, 36, 52, 21, 357},
	{1, 53, 59, 19, 133},
	{1, 60, 65, 18, 108},
	{1, 66, 70, 17, 85},
}

type D71 struct {
	disk *disk.Disk
	cbm  *Directory
}

func (d *D71) Read(reader *storage.Reader) error {
	if err := d.disk.Read(reader); err != nil {
		return fmt.Errorf("error reading the disk: %w", err)
	}

	// d.cbm = newDirectory(d.disk)
	d.cbm = &Directory{}
	if err := d.cbm.Read(); err != nil {
		return fmt.Errorf("error reading the directory: %w", err)
	}

	return nil
}
