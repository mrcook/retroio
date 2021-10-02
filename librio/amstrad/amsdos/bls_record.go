package amsdos

// BLS Table
//
// The values of BSH and BLM determine (implicitly) the data allocation
// size BLS, which is not an entry in the DPB.
type BlsRecord map[uint16]struct {
	Dirs uint16 // Reserved directory blocks
	BSH  uint8  // Block shift factor
	BLM  uint8  // Block mask
}

// TODO: where should this map be defined?
var BlsTable = BlsRecord{
	1024:  {Dirs: 32, BSH: 3, BLM: 7},    //
	2048:  {Dirs: 64, BSH: 4, BLM: 15},   // 0x0F
	4096:  {Dirs: 128, BSH: 5, BLM: 31},  // 0x1F
	8192:  {Dirs: 256, BSH: 6, BLM: 63},  // 0x3F
	16384: {Dirs: 512, BSH: 7, BLM: 127}, // 0x7F
}
