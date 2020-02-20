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
	disk *disk.Disk
	cbm  *Directory
}

func New(reader *storage.Reader, diskSize uint32) (*D81, error) {
	dsk, err := disk.New(reader, commodore.D81, diskSize)
	if err != nil {
		return nil, err
	}

	return &D81{disk: dsk}, nil
}

func (d *D81) Read() error {
	if err := d.disk.Read(); err != nil {
		return fmt.Errorf("error reading the disk: %w", err)
	}

	d.cbm = newDirectory(d.disk)
	if err := d.cbm.Read(); err != nil {
		return fmt.Errorf("error reading the directory: %w", err)
	}

	return nil
}

func (d D81) DisplayGeometry() {
	totalSectorCount := 0
	for _, t := range d.disk.Tracks {
		totalSectorCount += len(t.Sectors)
	}

	fmt.Println("DISK INFORMATION:")
	fmt.Println()
	fmt.Printf("Type:        %s\n", d.disk.DiskType())
	fmt.Printf("DOS Type:    %s\n", d.cbm.header.DosTypeDescription())
	fmt.Printf("Size:        %.2fKB\n", d.disk.DiskSizeInKB())
	fmt.Printf("Tracks:      %d\n", len(d.disk.Tracks))
	fmt.Printf("Sectors:     %d\n", totalSectorCount)
	fmt.Println()
	fmt.Printf("Name:        %s\n", d.cbm.header.PrintableDiskName())
	fmt.Printf("Files:       %d\n", len(d.cbm.directories))
	fmt.Printf("Free Blocks: %d\n", d.freeBlocks())
	fmt.Println()
}

func (d D81) CommandDir() {
	fmt.Println("LOAD\"$\",8")
	fmt.Println("SEARCHING FOR $")
	fmt.Println("LOADING")
	fmt.Println("READY.")
	fmt.Println("LIST")
	fmt.Println()
	fmt.Printf("0 \"%-16s\" %s %c%c\n", d.cbm.header.PrintableDiskName(), d.cbm.header.DiskID, d.cbm.header.DosVersion, d.cbm.header.DiskVersion)

	for _, dir := range d.cbm.directories {
		filename := fmt.Sprintf("\"%s\"", dir.Filename)
		fmt.Printf("%-3d  %-18s %s\n", dir.DirEntry.FileSizeInSectors, filename, dir.FileType)
	}

	fmt.Printf("%d BLOCKS FREE.\n", d.freeBlocks())
	fmt.Println()
}

func (d D81) freeBlocks() int {
	freeSectors := 0
	for _, b := range d.cbm.bamSide1.Entries {
		freeSectors += int(b.FreeSectors)
	}
	for _, b := range d.cbm.bamSide2.Entries {
		freeSectors += int(b.FreeSectors)
	}
	return freeSectors
}
