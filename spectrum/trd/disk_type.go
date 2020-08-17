package trd

// DiskType represents the possible values of the disk type
type DiskType byte

const (
	Tracks80Sides2 DiskType = 0x16
	Tracks40Sides2 DiskType = 0x17
	Tracks80Sides1 DiskType = 0x18
	Tracks40Sides1 DiskType = 0x19
)

func (dt DiskType) String() string {
	switch dt {
	case Tracks80Sides2:
		return "80 tracks, double side"
	case Tracks40Sides2:
		return "40 tracks, double side"
	case Tracks80Sides1:
		return "80 tracks, single side"
	case Tracks40Sides1:
		return "40 tracks, single side"
	}

	return "unknown"
}
