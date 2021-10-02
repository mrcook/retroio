package disk

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
