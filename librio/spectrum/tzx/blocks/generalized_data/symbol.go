package generalized_data

// The alphabet is stored using a table where each symbol is a row of pulses. The number of columns
// (i.e. pulses) of the table is the length of the longest sequence amongst all (MAXP=NPP or NPD,
// for pilot/sync or data blocks respectively); shorter waves are terminated by a zero-length pulse
// in the sequence.
// Any number of data symbols is allowed, so we can have more than two distinct waves; for
// example, imagine a loader which writes two bits at a time by encoding them with four distinct
// pulse lengths: this loader would have an alphabet of four symbols, each associated to a specific
// sequence of pulses (wave).
type Symbol struct {
	// Symbol flags:
	//   b0-b1: starting symbol polarity
	//   00:    opposite to the current level (make an edge, as usual) - default
	//   01:    same as the current level (no edge - prolongs the previous pulse)
	//   10:    force low level
	//   11:    force high level
	Flags uint8

	PulseLengths []uint16 // WORD[MAXP] Array of pulse lengths.
}
