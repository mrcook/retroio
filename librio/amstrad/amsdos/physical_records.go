package amsdos

// Map for SectorSize to PSH / PHM values
//
// PSH specifies the physical record shift factor: LOG2(sector-size/128).
// PHM specifies the physical record mask:         (sector-size/128)-1.
type PhysicalRecords map[uint16]struct {
	PSH uint8
	PHM uint8
}

var PhysicalShiftMaskTable = PhysicalRecords{
	128:  {PSH: 0, PHM: 0},
	256:  {PSH: 1, PHM: 1},
	512:  {PSH: 2, PHM: 3},
	1024: {PSH: 3, PHM: 7},
	2048: {PSH: 4, PHM: 15},
	4096: {PSH: 5, PHM: 31},
}
