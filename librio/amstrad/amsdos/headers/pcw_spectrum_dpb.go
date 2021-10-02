package headers

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
