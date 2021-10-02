package custom_info

// Picture block
//
// Finally you can include any picture (in supported formats) in the TZX file too.
// So cover pictures, maps, etc. can now be included in full colour (or whatever
// the formats supports). The best way for utilities to use this block is to spawn
// an external viewer, or the authors can write their own viewers (yeah, right ;).
// For inlay cards and other pictures that have zillions of colours use the
// JPEG format, for more simple pictures (drawing, maps, etc.) use the GIF format.
type Picture struct {
	// `Picture` + 9 spaces (custom block ID)
	Signature [16]byte

	DataLength        uint32  // Length of data that follow
	Format            uint8   // Picture format: 00 - GIF 01 - JPEG
	DescriptionLength uint8   // Description length (if this is 0 then handle it as 'Inlay Card')
	Description       []byte  // Description of the document in ASCII format
	Data              []uint8 // The picture itself
}
