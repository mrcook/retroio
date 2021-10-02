package tap

import "github.com/mrcook/retroio/librio/storage"

// Block is an interface for TAP header/data block
type Block interface {
	Read(reader *storage.Reader)
	Id() uint8
	Filename() string
	Name() string
	BlockData() []byte
}

// A Block as stored on tape may be a header or any data from the ZX Spectrum.
// NOTE: these blocks are also re-used in the TZX tape format.
//
// Blocks must have at least 2 bytes to build up a valid tape data block
// as each block should be bookended by `flag` and `checksum` bytes.
// Shorter data blocks (0 or 1 byte long) are so-called _fragmented_ ones.
type TapeBlock struct {
	Length   uint16
	TapeData Block // TODO: needs to be a generic type
}
