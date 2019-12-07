package spectrum

type Image interface {
	Read() error
	DisplayImageMetadata()
}
