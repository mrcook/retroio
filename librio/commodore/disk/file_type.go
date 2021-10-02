package disk

type FileType struct {
	Value       uint8
	Type        string
	SaveFlag    bool
	LockedFlag  bool
	ClosedFlag  bool
	Description string
}
