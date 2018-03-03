// Package contains TZX header data block
package tape

// Header is the first block of data found in all TZX files.
// The file is identified with the first 8 bytes being 'ZXTape!' plus the 'end of file' byte 26 (1A hex).
// This is followed by two bytes containing the major and minor version numbers.
type Header struct {
	Signature    [7]byte
	Terminator   uint8 // End of file marker
	MajorVersion uint8
	MinorVersion uint8
}
