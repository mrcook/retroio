// Package tap implements reading of Commodore C64 tape files as specified at:
// https://ist.uwaterloo.ca/~schepers/formats/TAP.TXT
//
// Designed by Per Hakan Sundell (author of the CCS64 C64 emulator) in 1997,
// this format attempts to duplicate the data stored on a C64 cassette tape,
// bit for bit. Since it is simply a representation of the raw serial data
// from a tape, it should handle *any* custom tape loaders that exist.
//
// The TAP images are generally very large, being a minimum of eight times,
// and up to sixteen times as large as what a raw PRG file would be. This is
// due to the way the data is stored, with each bit of the original file now
// being one byte large in the TAP file. The layout is fairly simple, with a
// small 14-byte header followed by file data.
//
// Note: DWORD values in the header are stored in little endian byte order.
package tap

// TAP file format
//
// The TAP format attempts to duplicate the data stored on a C64 cassette tapes
// bit for bit. Since it is simply a representation of the raw serial data from
// a tape, it should handle *any* custom tape loaders that exist.
//
// TAP images are generally very large, being a minimum of eight times, and up
// to sixteen times as large as what a raw PRG file would be. This is due to
// the way the data is stored, with each bit of the original file now being one
// byte large in the TAP file. The layout is fairly simple, with a small
// 20-byte header followed by file data.
type TAP struct {
	//
	// TAP Header (20-bytes)
	//

	// Signature: `C64-TAPE-RAW`
	Signature [12]byte

	Version    uint8    // TAP version: $00: original layout, $01: updated layout.
	Reserved   [3]uint8 // Reserved for future expansion
	DataLength uint32   // File data size (not including the header)

	//
	// TAP file data
	//
	Data []uint8
}
