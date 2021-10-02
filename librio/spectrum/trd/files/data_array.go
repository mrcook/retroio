package files

// File body of data array files.
type DataArray struct {
	Data         []uint8 // bytes representing the content of the data array
	Indicator    []uint8 // CHAR[] parameter 2 indicator; always #128/#170
	Unused       []uint8 // unused
	VariableName []uint8 // bits 0..5:  1..26 meaning "a".."z"
}
