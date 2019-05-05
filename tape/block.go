package tape

// Block is an interface for TZX data blocks
type Block interface {
	Id() uint8
	Name() string
	ToString() string
}
