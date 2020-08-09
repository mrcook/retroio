// Package disk implements reading of Commodore D64, D71, and D81 disk images.
//
// The D64, D71, D81 disk formats were documented by Peter Schepers. These
// specifications can be found at:
// http://vice-emu.sourceforge.net/vice_16.html
// https://ist.uwaterloo.ca/~schepers/formats.html
package disk

import (
	"errors"
	"fmt"

	"github.com/mrcook/retroio/commodore"
	"github.com/mrcook/retroio/storage"
)

// Disk object for working with the Commodore disk storage media.
type Disk struct {
	reader *storage.Reader

	Variation layout
	Tracks    []Track
}

// Initializes a new disk using the reader and media type/size values.
// The file size is required as each media type has at least two variations:
// e.g. with and without error bytes.
func New(reader *storage.Reader, mediaType commodore.MediaType, fileSize uint32) (*Disk, error) {
	layout, err := layoutForMedia(mediaType, fileSize)
	if err != nil {
		return nil, errors.New("invalid file - unexpected file size")
	}

	d := &Disk{
		reader:    reader,
		Variation: layout,
	}

	return d, nil
}

// Read the disk image, returning errors if the disk layout (track/sector counts)
// are of unexpected values for the expected media type.
func (d *Disk) Read() error {
	var i uint8

	// NOTE: track counting starts at 1, not 0!
	for i = 1; i <= d.Variation.tracks; i++ {
		geometry, err := trackLayout(d.Variation.mediaType, i)
		if err != nil {
			return fmt.Errorf("unable to determine geomery for this track number: %w", err)
		}
		track := Track{Number: i, Geometry: geometry}

		if err := track.Read(d.reader); err != nil {
			return fmt.Errorf("error reading track #%d: %w", i, err)
		}

		d.Tracks = append(d.Tracks, track)
	}

	totalSectorCount := 0
	for _, track := range d.Tracks {
		totalSectorCount += len(track.Sectors)
	}
	if totalSectorCount != int(d.Variation.totalSectors) {
		return fmt.Errorf("incorrect sector count for disk, expected %d, got %d\n", d.Variation.totalSectors, totalSectorCount)
	}

	return nil
}

func (d Disk) DiskType() string {
	return d.Variation.description
}

func (d Disk) DiskSizeInKB() float64 {
	return float64(d.Variation.diskSize) / 1024
}
