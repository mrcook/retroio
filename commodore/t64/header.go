package t64

import (
	"encoding/binary"
	"fmt"

	"github.com/mrcook/retroio/storage"
)

// Tape Header
type Header struct {
	Signature   [32]byte // DOS tape description + EOF (for type)
	Version     uint16   // Tape version (always $0200 for VICE)
	MaxEntries  uint16   // Maximum number of entries in the directory
	UsedEntries uint16   // Total used entries. Used = $0005 = 5 entries.
	Unused      uint16   // Not Used
	Name        [24]byte // Tape container name, padded with $20 (space)
}

func (h *Header) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, h)
}

func (h Header) String() string {
	str := ""
	str += fmt.Sprintf("Name:            %s\n", h.Name)
	str += fmt.Sprintf("Version:         $%04x\n", h.Version)
	str += fmt.Sprintf("Max Directories: %d\n", h.MaxEntries)
	str += fmt.Sprintf("Used Entries:    %d\n", h.UsedEntries)
	return str
}
