package amsdos

// Default DPB values for the Amstrad CPC SSSD disk format.
const (
	ExtentMask      uint8  = 0
	TrackCount      uint8  = 40
	SectorsPerTrack uint8  = 9
	BytesPerSector  uint16 = 512
	ReadWriteGap    uint8  = 0x2A // uPD765A read/write gap
	FormatGap       uint8  = 0x52 // uPD765A format gap

	DiskByteSize = int(TrackCount) * int(SectorsPerTrack) * int(BytesPerSector)

	BLS uint16 = 1024
	DSM uint16 = uint16(DiskByteSize / int(BLS))
	DRM uint16 = 64 // All formats have 64 directory entries: 64 * 32-bytes = 2048 bytes.
)

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
	Name          [8]uint8  // Name part, padded with spaces
	Type          [3]uint8  // Type part, padded with spaces
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

// When a file without a header is opened for input a fake header is constructed in store.
// TODO: probably not needed, just use the normal disc header
type HeaderlessHeader struct {
	// Filename
	User    uint8    // User number, #00..#0F
	Name    [8]uint8 // Name part, padded with spaces
	Type    [3]uint8 // Type part, padded with spaces
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

// PCW/Spectrum system
//
// In addition to the XDPB system, the PCW and Spectrum +3 can determine the format
// of a disc from a 16-byte record on track 0, head 0, physical sector 1.
//
// If all bytes of the spec are 0E5h, it should be assumed that the disc is a
// 173k PCW/Spectrum +3 disc, ie:
//   single sided, single track, 40 tracks, 9 sectors/track, 512-byte sectors,
//   1 reserved track, 1k blocks,
//   2 directory blocks,
//   gap lengths 2Ah and 52h,
//   not bootable
//
// PCW16 extended boot record
//
// The "boot record" system has been extended in PCW16 CP/M (BIOS 0.09 and later).
// The extension is intended to allow a CP/M "partition" on a DOS-formatted floppy disc.
//
// An extended boot sector (cylinder 0, head 0, sector 1) has the following characteristics:
// - First byte is 0E9h or 0EBh
// - Where DOS expects the disc label to be (at sector + 2Bh) there are 11 ASCII bytes
//   of the form `CP/M????DSK`, where "?" can be any character.
// - At sector + 7Ch are the four ASCII bytes "CP/M"
// - At sector + 80h is the disc specification as described above.
type PcwSpectrumDPB struct {
	// format number
	//   0: SS SD
	//   1: CPC formats, but those formats don't have boot records anyway.
	//   2: ^
	//   3: DS DD
	// Any other value: bad format
	FormatNumber uint8

	// sidedness ; As in XDPB
	MediaType uint8

	// tracks/side
	TrackCountPerSide uint8

	// sectors/track
	SectorCountPerTrack uint8

	// physical sector shift ; psh in XDPB
	PhysicalShift uint8

	// no. reserved tracks ; off in XDPB
	ReservedTracks uint8

	// block shift ; bsh in XDPB
	BlockShift uint8

	// no. directory blocks
	DirectoryBlockCount uint8

	// uPD765A read/write gap length
	ReadWriteGap uint8

	// uPD765A format gap length
	FormatGap uint8

	// 0,0,0,0,0 ; Unused
	Unused [5]uint8

	// Checksum fiddle byte ; Used to indicate Bootable discs.
	//
	// Change this byte so that the 8-bit checksum of the sector is:
	//    1 - sector contains a PCW9512 bootstrap
	//    3 - sector contains a Spectrum +3 bootstrap
	//  255 - sector contains a PCW8256 bootstrap
	//        (the bootstrap code is in the remainder of the sector)
	Checksum uint8
}
