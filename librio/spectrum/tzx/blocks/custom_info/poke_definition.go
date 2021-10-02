package custom_info

// POKE block: POKE definition (0x05 bytes)
//
// For each POKE entry you can supply the memory page number and/or the original
// value of the address (if you want to restore it some way through the game).
// Normally you would enter these pokes with the help of some freezer-type tool
// like Multiface, but hopefully in the future the emulators will support this
// block directly, in which case you could use the 'User inserts the POKE value'
// feature. You can specify the point at which to insert the POKEs in the 'General
// description' field.
type PokeDefinition struct {
	// POKE type:
	//   bit 0-2 : memory page number
	//   bit 3 : ignore memory page number
	//   bit 4 : user inserts the POKE value
	//   bit 5 : unknown original value
	PokeType uint8

	Address       uint16 // POKE address
	Value         uint8  // POKE value (leave 0 if 'user inserts' bit set)
	OriginalValue uint8  // POKE original value (leave 0 if 'unknown' bit set)
}
