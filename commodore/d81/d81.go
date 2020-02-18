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

	"retroio/commodore"
	"retroio/commodore/disk"
	"retroio/storage"
)

type D81 struct {
	Disk *disk.Disk
}

func New(reader *storage.Reader, diskSize uint32) (*D81, error) {
	dsk, err := disk.New(reader, commodore.D81, diskSize)
	if err != nil {
		return nil, err
	}

	return &D81{Disk: dsk}, nil
}

func (d *D81) Read() error {
	return d.Disk.Read()
}

func (d D81) DisplayGeometry() {
	d.Disk.DisplayGeometry()
}

func (d D81) CommandDir() {
	fmt.Println("directory listing unsupported")
}
