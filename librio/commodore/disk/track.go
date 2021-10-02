package disk

import (
	"fmt"

	"github.com/mrcook/retroio/librio/storage"
)

type Track struct {
	Sectors []Sector

	// TODO: are these part of the specification?
	Number   uint8
	Geometry TrackGeometry
}

func (t *Track) Read(reader *storage.Reader) error {
	for i := 0; i < t.Geometry.SectorsPerTrack; i++ {
		sector := Sector{}

		if err := sector.Read(reader); err != nil {
			return fmt.Errorf("error reading sector #%d: %w", i, err)
		}

		t.Sectors = append(t.Sectors, sector)
	}

	return nil
}
