package disk

// A BAM Entry for each track on the disk
//
// The sector bitmap bytes (three for the D64/D71 and five for D81) indicate
// which sectors are used/free for one track.
//
// A D81 disk with a five byte bitmap gives a total of 40 bits of storage
// (8 bits per byte), one for each sector on the track. If track 40 has these
// five bytes set as $F0, $FF, $2D, $FF, $FE, they would produce a bitmap of:
//
//   F0=11110000, FF=11111111, 2D=00101101, FF=11111111, FE=11111110
//
// The sector bits of *each byte* are stored right-to-left, so to aid in
// understanding the binary notation, flip the bits around.
//
//              111111 11112222 22222233 33333333
//   01234567 89012345 67890123 45678901 23456789
//   -------------------------- -------- --------
//   00001111 11111111 10110100 11111111 01111111
//   ^                     ^             ^
//   sector 0           sector 20     sector 32
//
// The 40 sectors of the first track use bit positions 0-39. If a bit is
// set (1), then the sector is free. Therefore, track 40 has sectors 0-3,
// 17, 20, 22, 23, and 32 in use, the rest are unused.
//
// A D64/D71 would use up to 21 bits from its three bitmap bytes (bits 0-20).
// The bits for non-existent sectors on a track are marked as allocated (reset).
//
// See `docs.md` for lots more information on this set of bytes.
type BamEntry24Bit struct {
	// The number of free sectors on this track
	FreeSectors uint8 // 0xff

	// Bitmap of which sectors are used/free.
	//   - D64 = 24-bit (3 bytes)
	//   - D71 = 24-bit (3 bytes)
	//   - D81 = 40-bit (5 bytes)
	SectorBitmap []uint8
}
