package custom_info

// Spectrum screen block
//
// If the game on the tape is not an original and lacks the original loading
// screen then you can supply it separately within this block. This is also very
// handy when you want the loading screen stored separately because the original
// is either encrypted (like with the 'Speedlock' or 'Alkatraz' loaders) or it is
// corrupted by some on-screen info (like the 'Bleepload' loader). Of course not
// only loading screens can be stored here... you can use it to store maps or any
// other picture that is in Spectrum Video format (that's why the Description is
// there for), but because the Loading Screen will be the most common you can
// just set the description length field to 0 when you use it for that. Also the
// border colour can be specified.
type ScreenBlock struct {
	// // `Spectrum Screen` + 1 spaces (custom block ID)
	Signature [16]byte

	DataLength        uint32      // Length of data that follow
	DescriptionLength uint8       // Description length (if this is 0 then handle it as 'Loading Screen')
	Description       []byte      // Description of the picture in ASCII format
	BorderColour      uint8       // BORDER Colour in Spectrum colour format (0=black, 1=blue, ...)
	Data              [6912]uint8 // Screen in standard Spectrum video format
}
