package headers

// AMDSDOS File Record Header
//
// Files may, or may not, have a header depending on the contents of the
// file - CP/M files do not have headers. This will not cause problems for
// programs written in BASIC but it is an important difference between
// cassette and disc files.
//
// AMSDOS files have a single header in the first 128 bytes of the file - the
// header record - except unprotected ASCII files, which have no header.
//
// These headers are detected by calculating the checksum the first 67 bytes of
// the record. If the checksum is as expected then a header is present, if not
// then there is no header. Thus it is possible, though unlikely, that a file
// without a header could be mistaken for one with a header.
type RecordHeader struct {
	// Cassette/Disc header
	User          uint8     // User number, #00..#0F
	Name          [8]byte   // Name part, padded with spaces
	Type          [3]byte   // Type part, padded with spaces
	Unknown       [4]uint8  // #00
	BlockNumber   uint8     // Not used, set to 0
	LastBlock     uint8     // Not used, set to 0
	FileType      uint8     // As per cassette
	DataLength    uint16    // As per cassette
	DataLocation  uint16    // As per cassette
	FirstBlock    uint8     // Set to #FF, only used for output files
	LogicalLength uint16    // As per cassette
	EntryAddress  uint16    // As per cassette
	Unallocated   [36]uint8 // As per cassette

	FileLength [3]uint8  // 24-bit value. Length of the file in bytes, excluding the header record. Least significant byte in lowest address.
	Checksum   uint16    // Sixteen bit checksum, sum of bytes 0..66
	Undefined  [58]uint8 // 69... 127 Undefined
}
