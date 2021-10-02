package custom_info

// POKEs Block
//
// The purpose of this custom block is to hold any amount of different trainers
// for the game.
//
// NOTE: All ASCII Descriptions can use more than one line. Please use only up to
// 30 characters per line and separate the lines by one CR (13dec).
//
// POKEs block: Trainer (0x15 bytes)
type Poke struct {
	// `POKEs` + 11 spaces (custom block ID)
	Signature [16]byte

	Length            uint32        // Length of data that follow
	DescriptionLength uint8         // General description length
	Description       []byte        // General description in ASCII format
	TrainerCount      uint8         // Number of trainers
	Trainers          []PokeTrainer // Trainer definitions
}
