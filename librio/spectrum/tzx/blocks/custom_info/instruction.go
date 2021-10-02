package custom_info

// Instructions block
//
// This block can hold any general .TXT file, with the main purpose of storing the
// instructions to the program or game that is in the tape. To ensure consistency
// with all other ASCII texts in this format please use a single `CR` character
// (`13` dec, `0D` hex) to separate lines; also please use only up to 80
// characters per line.
type Instruction struct {
	// `Instructions` + 4 spaces (custom block ID)
	Signature [16]byte

	Length       uint32 // Length of the following data
	Instructions string // Instructions text in ASCII format
}
