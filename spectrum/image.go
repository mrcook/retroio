package spectrum

type Image interface {
	Read() error
	DisplayGeometry()
	DisplayBASIC()
}
