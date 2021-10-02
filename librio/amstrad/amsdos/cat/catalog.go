package cat

type Catalog struct {
	Drive       byte
	User        uint8
	FreeSpace   uint16
	HiddenFiles int
	Records     []DirectoryRecord
}
