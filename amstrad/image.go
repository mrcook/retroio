package amstrad

type Image interface {
	Read() error
	DisplayImageMetadata()
}
