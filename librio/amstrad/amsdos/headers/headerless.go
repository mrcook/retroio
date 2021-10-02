package headers

// TODO: probably not needed, just use the normal disc header
// When a file without a header is opened for input a fake header is constructed in store.
type HeaderlessHeader struct {
	// Filename
	User    uint8    // User number, #00..#0F
	Name    [8]byte  // Name part, padded with spaces
	Type    [3]byte  // Type part, padded with spaces
	Unknown [4]uint8 // #00

	Unused1      uint8 // Not used, set to 0
	Unused2      uint8 // Not used, set to 0
	FileType     uint8 // #16, unprotected ASCII version 1
	Unused3      uint16
	DataLocation uint16 // Address of 2K buffer
	FirstBlock   uint8  // #FF
	Unused4      uint16
	Unused5      uint16
	Unused6      [36]uint8
}
