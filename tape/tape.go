package tape

type Tape interface {
	ReadBlocks() error
	DisplayTapeMetadata()
}
