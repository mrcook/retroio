// Package t64 implements reading of Commodore C64 tape files as specified at:
// http://vice-emu.sourceforge.net/vice_16.html#SEC295
// https://ist.uwaterloo.ca/~schepers/formats/T64.TXT
//
// The T64 File Structure was developed by Miha Peternel for use in the C64S emulator.
// It is easy to use and allows future extensions.
//
// Note: all WORD and DWORD values are stored in low/high byte order.
package t64

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/mrcook/retroio/librio/storage"
)

// T64 File structure
type T64 struct {
	// Tape Header
	Signature   [32]byte // DOS tape description + EOF (for type)
	Version     uint16   // Tape version (always $0200 for VICE)
	MaxEntries  uint16   // Maximum number of entries in the directory
	UsedEntries uint16   // Total used entries. Used = $0005 = 5 entries.
	Unused      uint16   // Not Used
	Name        [24]byte // Tape container name, padded with $20 (space)

	// File records for 32*n directory entries
	Records []Record

	// Binary data for the records
	Data [][]byte
}

func (t *T64) Read(reader *storage.Reader) error {
	for i := range t.Signature {
		t.Signature[i] = reader.ReadByte()
	}
	t.Version = reader.ReadShort()
	t.MaxEntries = reader.ReadShort()
	t.UsedEntries = reader.ReadShort()
	t.Unused = reader.ReadShort()
	for i := range t.Name {
		t.Name[i] = reader.ReadByte()
	}

	// Read headers for the records
	for i := 0; i < int(t.UsedEntries); i++ {
		r := Record{}
		if err := r.Read(reader); err != nil {
			return fmt.Errorf("binary.Read failed: %v", err)
		}
		t.Records = append(t.Records, r)
	}

	// In case the records were stored out of order, sort 'em.
	sort.Slice(t.Records, func(i, j int) bool {
		return t.Records[i].Offset < t.Records[j].Offset
	})

	if err := t.readDataEntries(reader); err != nil {
		return err
	}

	return nil
}

// readDataEntries reads the data for each record.
// TODO: improve this crufty code
func (t *T64) readDataEntries(reader *storage.Reader) error {
	// Add a small failsafe before continuing
	if len(t.Records) == 0 {
		return fmt.Errorf("can not read data, no records available")
	}

	// Get the current offset value
	headerLength := 64
	recordLength := 32 * len(t.Records)
	offset := headerLength + recordLength

	for _, r := range t.Records {
		// Discard any bytes before the start of the record
		// First record typically starts at $0400, but not always!
		discardCount := int(r.Offset) - offset
		if discardCount > 0 {
			if _, err := reader.Discard(discardCount); err != nil {
				return err
			}
		}

		data, err := r.DataRead(reader, offset)
		if err != nil {
			return errors.Wrap(err, "error reading record data")
		}

		// set offset to end of current record
		offset = int(r.Offset) + len(data)

		t.Data = append(t.Data, data)
	}

	return nil
}
