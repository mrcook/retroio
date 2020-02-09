package disk

import (
	"errors"
	"fmt"

	"retroio/storage"
)

// Disk object for working with the Commodore disk format
// Can work with the D64, D71, D81 images.
type Disk struct {
	Reader    *storage.Reader
	Variation layout
	Tracks    []Track
}

func New(fileSize uint32, diskType storageID, reader *storage.Reader) (*Disk, error) {
	variation, err := layoutByFileSize(fileSize)
	if err != nil {
		return nil, errors.New("invalid file; unexpected file size")
	}

	if variation.format != diskType {
		return nil, fmt.Errorf("invalid file size for a %d", diskType)
	}

	d := &Disk{
		Reader:    reader,
		Variation: variation,
	}
	return d, nil
}

func (d *Disk) Read() error {
	var i uint8

	// Track counting starts at 1 (not 0).
	for i = 1; i <= d.Variation.trackCount; i++ {
		trackGeometry, err := trackLayout(d.Variation.format, i)
		if err != nil {
			return fmt.Errorf("no geomery info found for this track number: %w", err)
		}

		track := Track{Number: i, Geometry: trackGeometry}
		if err := track.Read(d.Reader); err != nil {
			return fmt.Errorf("error reading track #%d: %w", i, err)
		}

		d.Tracks = append(d.Tracks, track)
	}

	totalSectorCount := 0
	for _, track := range d.Tracks {
		totalSectorCount += len(track.Sectors)
	}
	if totalSectorCount != int(d.Variation.totalSectorCount) {
		return fmt.Errorf("incorrect sector count, expected %d, got %d\n", d.Variation.totalSectorCount, totalSectorCount)
	}

	return nil
}

// DisplayGeometry prints the disk, track and sector metadata to the terminal.
func (d Disk) DisplayGeometry() {
	totalSectorCount := 0
	for _, t := range d.Tracks {
		totalSectorCount += len(t.Sectors)
	}

	fmt.Println("DISK INFORMATION:")
	fmt.Println()
	fmt.Printf("Type:      %s\n", d.Variation.description)
	fmt.Printf("Disk Size: %.2fKB\n", float64(d.Variation.fileSize)/1024)
	fmt.Printf("Tracks:    %d\n", len(d.Tracks))
	fmt.Printf("Sectors:   %d\n", totalSectorCount)
}
