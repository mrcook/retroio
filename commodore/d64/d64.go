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

	"github.com/mrcook/retroio/commodore"
	"github.com/mrcook/retroio/commodore/disk"
	"github.com/mrcook/retroio/storage"
)

type D64 struct {
	disk *disk.Disk
	cbm  *Directory
}

func New(reader *storage.Reader, diskSize uint32) (*D64, error) {
	dsk, err := disk.New(reader, commodore.D64, diskSize)
	if err != nil {
		return nil, err
	}

	return &D64{disk: dsk}, nil
}

func (d *D64) Read() error {
	if err := d.disk.Read(); err != nil {
		return fmt.Errorf("error reading the disk: %w", err)
	}

	d.cbm = newDirectory(d.disk)
	if err := d.cbm.Read(); err != nil {
		return fmt.Errorf("error reading the directory: %w", err)
	}

	return nil
}

func (d D64) DisplayGeometry() {
	totalSectorCount := 0
	for _, t := range d.disk.Tracks {
		totalSectorCount += len(t.Sectors)
	}

	fmt.Println("DISK INFORMATION:")
	fmt.Println()
	fmt.Printf("Type:        %s\n", d.disk.DiskType())
	fmt.Printf("DOS Type:    %s\n", d.cbm.bam.DosTypeDescription())
	fmt.Printf("Size:        %.2fKB\n", d.disk.DiskSizeInKB())
	fmt.Printf("Tracks:      %d\n", len(d.disk.Tracks))
	fmt.Printf("Sectors:     %d\n", totalSectorCount)
	fmt.Println()
	fmt.Printf("Name:        %s\n", d.cbm.bam.PrintableDiskName())
	fmt.Printf("Files:       %d\n", len(d.cbm.directories))
	fmt.Printf("Free Blocks: %d\n", d.freeBlocks())
	fmt.Println()
}

func (d D64) CommandDir() {
	fmt.Println("LOAD\"$\",8")
	fmt.Println("SEARCHING FOR $")
	fmt.Println("LOADING")
	fmt.Println("READY.")
	fmt.Println("LIST")
	fmt.Println()
	fmt.Printf("0 \"%-16s\" %s %c%c\n", d.cbm.bam.PrintableDiskName(), d.cbm.bam.DiskID, d.cbm.bam.DosVersion, d.cbm.bam.DiskVersion)

	for _, dir := range d.cbm.directories {
		filename := fmt.Sprintf("\"%s\"", dir.Filename)
		fmt.Printf("%-3d  %-18s %s\n", dir.DirEntry.FileSizeInSectors, filename, dir.FileType)
	}

	fmt.Printf("%d BLOCKS FREE.\n", d.freeBlocks())
	fmt.Println()
}

func (d D64) freeBlocks() int {
	freeSectors := 0
	for i, b := range d.cbm.bam.Entries {
		if i == DirectoryTrackNumber {
			continue
		}
		freeSectors += int(b.FreeSectors)
	}
	return freeSectors
}
