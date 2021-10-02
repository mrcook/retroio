package disk

// SideSector of a REL file type.
//
// Suitable for D64, D71, and D81 formats.
//
// For each REL file there are up to 6 SIDE SECTORS / SUPER SIDE SECTORS.
// Each side-sector has 120 possible Track/Sector chains, allowing for a
// maximum REL file size of 720 data sectors.
type SideSector struct {
	// Track location of next side-sector ($00 if last sector)
	TrackLocation uint8

	// Sector location of next side-sector
	SectorLocation uint8

	// Side-sector block number
	// first sector is $00, the next is $01, $02, etc.
	BlockNumber uint8

	// REL file RECORD size (from directory entry, max. value 254)
	RecordSize uint8

	// Track/sector locations of all six side-sectors.
	// Note the first entry is this very sector we have listed here.
	// The next is the next t/s listed at the beginning of the sector.
	// All of this information must be correct.
	// If one of these chains is $00/$00, then we have no more side sectors.
	// Also, all of these (up to six) side sectors must have the same values
	// in this range.
	AllSideSectorLocations [6][2]uint8

	// T/S chains of *each* sector of the data portion.
	// When we get a $00/$00, we are at the end of the file.
	TrackSectorChains [120][2]uint8
}
