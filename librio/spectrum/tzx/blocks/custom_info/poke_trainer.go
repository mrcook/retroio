package custom_info

// Each trainer can have its own description and any number of POKEs.
type PokeTrainer struct {
	DescriptionLength uint8            // Trainer description length
	Description       []byte           // Trainer description in ASCII format
	PokeCount         uint8            // Number of pokes in this trainer
	PokeDefinitions   []PokeDefinition // POKEs definitions
}
