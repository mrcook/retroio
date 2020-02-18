package commodore

// Custom type for the different storage mediums
type MediaType uint8

const (
	Unknown MediaType = iota
	T64
	TAP
	D64 // all variations: 35, 40, and 42 tracks
	D71
	D81
)

type Image interface {
	Read() error
	DisplayGeometry()
	CommandDir()
}
