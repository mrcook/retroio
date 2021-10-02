package files

// File body of a BASIC program.
type BasicProgram struct {
	Data                []uint8  // bytes of the BASIC program (and if present, variables)
	Indicator           []uint8  // CHAR[] parameter 2 indicator; always #128/#170
	AutostartLineNumber []uint16 // 0..9999
}
