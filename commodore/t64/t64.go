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

	"github.com/mrcook/retroio/storage"
)

// T64 File structure
type T64 struct {
	reader *storage.Reader

	Header  Header   // The Tape header
	Records []Record // File records for 32*n directory entries
	Data    [][]byte // Binary data for the records

}

func New(reader *storage.Reader) *T64 {
	return &T64{reader: reader}
}

func (t *T64) Read() error {
	t.Header = Header{}
	if err := t.Header.Read(t.reader); err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	// Read headers for the records
	for i := 0; i < int(t.Header.UsedEntries); i++ {
		r := Record{}
		if err := r.Read(t.reader); err != nil {
			return fmt.Errorf("binary.Read failed: %v", err)
		}
		t.Records = append(t.Records, r)
	}

	// In case the records were stored out of order, sort 'em.
	sort.Slice(t.Records, func(i, j int) bool {
		return t.Records[i].Offset < t.Records[j].Offset
	})

	if err := t.readDataEntries(); err != nil {
		return err
	}

	return nil
}

// DisplayGeometry prints the tape metadata and record headers to the terminal.
func (t T64) DisplayGeometry() {
	fmt.Println("HEADER INFORMATION:")
	fmt.Println(t.Header)

	for i, r := range t.Records {
		fmt.Printf("RECORD #%d:\n", i)
		fmt.Println(r)
	}

	for i, r := range t.Data {
		fmt.Printf("BINARY DATA #%d = %d bytes\n", i, len(r))
	}
}

func (t T64) CommandDir() {
	fmt.Println("directory listing unsupported")
}

// readDataEntries reads the data for each record.
// TODO: improve this crufty code
func (t *T64) readDataEntries() error {
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
			if _, err := t.reader.Discard(discardCount); err != nil {
				return err
			}
		}

		data, err := r.DataRead(t.reader, offset)
		if err != nil {
			return errors.Wrap(err, "error reading record data")
		}

		// set offset to end of current record
		offset = int(r.Offset) + len(data)

		t.Data = append(t.Data, data)
	}

	return nil
}
