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
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"retroio/storage"
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
	if err := t.readHeader(); err != nil {
		return err
	}
	if err := t.readRecords(); err != nil {
		return err
	}
	if err := t.readDataEntries(); err != nil {
		return err
	}

	return nil
}

// DisplayImageMetadata prints the tape metadata and record headers to the terminal.
func (t T64) DisplayImageMetadata() {
	fmt.Println("HEADER INFORMATION:")
	fmt.Println(t.Header.ToString())

	for i, r := range t.Records {
		fmt.Printf("RECORD #%d:\n", i)
		fmt.Println(r.ToString())
	}

	for i, r := range t.Data {
		fmt.Printf("BINARY DATA #%d = %d bytes\n", i, len(r))
	}
}

// readHeader reads the cassette tape header data.
func (t *T64) readHeader() error {
	t.Header = Header{}
	if err := binary.Read(t.reader, binary.LittleEndian, &t.Header); err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}
	return nil
}

// readRecords reads the record header entries.
func (t *T64) readRecords() error {
	for i := 0; i < int(t.Header.UsedEntries); i++ {
		r := Record{}
		if err := binary.Read(t.reader, binary.LittleEndian, &r); err != nil {
			return fmt.Errorf("binary.Read failed: %v", err)
		}
		t.Records = append(t.Records, r)
	}

	// In case the records were stored out of order, sort 'em.
	sort.Slice(t.Records, func(i, j int) bool {
		return t.Records[i].Offset < t.Records[j].Offset
	})

	return nil
}

// readDataEntries reads the data for each record.
// TODO: oh boy, this not pretty, but it works for now.
func (t *T64) readDataEntries() error {
	// Add a small failsafe before continuing
	if len(t.Records) == 0 {
		return fmt.Errorf("can not read data, no records available")
	}

	// Get the current offset value
	headerLength := 64
	recordLength := 32 * len(t.Records)
	currentOffset := headerLength + recordLength

	for _, r := range t.Records {
		// Discard any bytes before the start of the record
		// First record typically starts at $0400, but not always!
		discardCount := int(r.Offset) - currentOffset
		if discardCount > 0 {
			if _, err := t.reader.Discard(discardCount); err != nil {
				return err
			}
		}

		// Now read the correct number of bytes
		length := int(r.EndAddress - r.StartAddress)
		data := make([]byte, length)
		if _, err := io.ReadFull(t.reader, data); err != nil {
			if err == io.EOF {
				return nil
			} else if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		t.Data = append(t.Data, data)

		// set offset to end of current record
		currentOffset = int(r.Offset) + length
	}

	return nil
}
