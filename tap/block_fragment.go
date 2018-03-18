package tap

// FragmentDataBlocks for storing data without flag or checksum byte
// case #6: fragment data blocks
// Fragmented data blocks cannot be produced by ROM saving routines - they must
// be produced with machine code programs. They have less than 2 bytes. In some
// games you really find zero length fragment data blocks.
type FragmentDataBlocks struct {
	Data []byte // BYTE[]   len is < 2   data block the essential data (may be empty)
}

// ToString returns a formatted string for the block
func (b FragmentDataBlocks) ToString() string {
	return "not implemented"
}
