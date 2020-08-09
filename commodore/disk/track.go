package disk

import (
	"fmt"

	"github.com/mrcook/retroio/storage"
)

type Track struct {
	Number   uint8
	Geometry trackGeometry
	Sectors  []Sector
}

func (t *Track) Read(reader *storage.Reader) error {
	for i := 0; i < t.Geometry.sectorsPerTrack; i++ {
		sector := Sector{}

		if err := sector.Read(reader); err != nil {
			return fmt.Errorf("error reading sector #%d: %w", i, err)
		}

		t.Sectors = append(t.Sectors, sector)
	}

	return nil
}

func (t Track) String() string {
	str := ""
	str += fmt.Sprintf("Track #%02d : %d sectors", t.Number, len(t.Sectors))
	return str
}
