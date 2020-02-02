package amstrad

type Image interface {
	Read() error
	DisplayGeometry()
	CommandDir()
}
