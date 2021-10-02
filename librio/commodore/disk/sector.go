package disk

import (
	"encoding/binary"

	"github.com/mrcook/retroio/librio/storage"
)

// Sector stores the data for each sector of a track.
// The first two bytes store the next Track/Sector locations followed by 254 bytes of data.
type Sector struct {
	Data [256]uint8
}

func (s *Sector) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, s)
}
