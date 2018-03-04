package tape

// Block is an interface for TZX data blocks
type Block interface {
	Id() int
	Name() string
	Metadata() string
}
