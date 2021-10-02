package custom_info

// ZX-Edit document block
//
// This block can hold files created with the new utility called ZX-Editor. This
// utility gives documents the look and feel of ZX-Spectrum and its documents can
// contain text, graphics (with Spectrum attributes), different type faces,
// colours, etc. Normally these files use extension `.ZED`. Also the description is
// added, in case you want to use it for something else than 'Instructions' - you
// can use it for MAPs, etc.
type ZxEditDocument struct {
	// `ZX-Edit document` (custom block ID)
	Signature [16]byte

	DataLength        uint32  // Length of data that follow
	DescriptionLength uint8   // Description length (if this is 0 then handle it as 'Instructions')
	Description       []byte  // Description of the document in ASCII format
	Data              []uint8 // The ZX-Editor document (.ZED file)
}
