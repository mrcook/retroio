package commodore

type Image interface {
	Read() error
	DisplayImageMetadata()
}
