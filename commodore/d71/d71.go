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

	"retroio/commodore"
	"retroio/commodore/d64"
	"retroio/commodore/disk"
	"retroio/storage"
)

type D71 struct {
	disk *disk.Disk
	cbm  *Directory
}

func New(reader *storage.Reader, diskSize uint32) (*D71, error) {
	dsk, err := disk.New(reader, commodore.D71, diskSize)
	if err != nil {
		return nil, err
	}

	return &D71{disk: dsk}, nil
}

func (d *D71) Read() error {
	if err := d.disk.Read(); err != nil {
		return fmt.Errorf("error reading the disk: %w", err)
	}

	d.cbm = newDirectory(d.disk)
	if err := d.cbm.Read(); err != nil {
		return fmt.Errorf("error reading the directory: %w", err)
	}

	return nil
}

func (d D71) DisplayGeometry() {
	d.disk.DisplayGeometry()
}

func (d D71) CommandDir() {
	fmt.Println("LOAD\"$\",8")
	fmt.Println("SEARCHING FOR $")
	fmt.Println("LOADING")
	fmt.Println("READY.")
	fmt.Println("LIST")
	fmt.Println()
	fmt.Printf("0 \"%-16s\" %s %c%c\n", d.cbm.bam.PrintableDiskName(), d.cbm.bam.DiskID, d.cbm.bam.DosVersion, d.cbm.bam.DiskVersion)

	for _, dir := range d.cbm.directories {
		fileType := dir.FileType.Abbreviation
		if dir.FileType.Value == 0 {
			fileType = fmt.Sprintf("%s (%s)", fileType, dir.FileType.Label)
		}
		filename := fmt.Sprintf("\"%s\"", dir.Filename)

		fmt.Printf("%-3d  %-18s %s\n", dir.DirEntry.FileSizeInSectors, filename, fileType)
	}

	fmt.Printf("%d BLOCKS FREE.\n", d.freeBlocks())
	fmt.Println()
}

func (d D71) freeBlocks() int {
	freeSectors := 0
	for i, b := range d.cbm.bam.Entries {
		if i == d64.DirectoryTrackNumber {
			continue
		}
		freeSectors += int(b.FreeSectors)
	}

	// Now count side 2 free sectors - for tracks 36-70
	for i, b := range d.cbm.bam.FreeSectorCount {
		if i == bamExtraTrackNumber {
			continue
		}
		freeSectors += int(b)
	}
	return freeSectors
}
