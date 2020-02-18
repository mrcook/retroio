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

// SuperSideSector of a D81 REL file type
//
// The SUPER SIDE SECTOR provides 126 groups of SIDE SECTORS, allowing for file
// sizes of up to (theoretically) 90720 sectors, or about 22.15 Megabytes.
// The layout of each SIDE SECTOR pointed to is the same as the D64/D71.
type SuperSideSector struct {
	// Track/sector of first side sector in Group 0
	TrackLocation  uint8
	SectorLocation uint8

	Unknown uint8 // Always $FE

	// T/S chains of *each* SIDE_SECTORS in Groups 0-125
	// NOTE: t/s of Group 0 listed above must also be included in this list.
	TrackSectorGroupChains [126][2]uint8

	Unused uint8 // Unused (likely $00)
}
