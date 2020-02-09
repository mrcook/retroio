package disk

import (
	"encoding/binary"

	"retroio/storage"
)

type Sector struct {
	NextPointer uint16
	Data        [254]byte
}

func (s *Sector) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, s)
}
