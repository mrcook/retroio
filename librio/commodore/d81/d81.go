// Package d81 implements reading Commodore D81 image files.
//
// The D81 is a sector-for-sector copy of a 1581 disk, with 80 tracks of 40
// sectors each (0 to 39) making a DOS file size of 819200 bytes. If the error
// byte block (3200 bytes) is attached, the file size wil be 822400 bytes.
//
// Additional D81 documentation can be found in the `docs.md` file.
package d81

import (
	"fmt"

	"github.com/mrcook/retroio/librio/commodore/disk"
	"github.com/mrcook/retroio/librio/storage"
)

const (
	DirectoryTrackNumber = 40 - 1 // tracks start counting from 1 not zero
	DirEntriesPerSector  = 8
	BamSectorBitmapSize  = 5 // 40-bit BAM bitmap
)

var Geometry = []disk.TrackGeometry{
	{0, 1, 40, 40, 1600},
	{1, 41, 80, 40, 1600},
}

type D81 struct {
	disk *disk.Disk
	cbm  *Directory
}

func (d *D81) Read(reader *storage.Reader) error {
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
