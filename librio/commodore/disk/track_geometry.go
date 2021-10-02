package disk

// TrackGeometry object to define the basic attributes of a range of tracks on disk:
//   - which side of the disk they are located
//   - sector information for the range
//   - total sector count for the range
//
// At present this only supports the Commodore disk formats.
type TrackGeometry struct {
	Side            uint8 // NOTE: starts at `0` not `1`
	StartTrack      uint8
	EndTrack        uint8
	SectorsPerTrack int
	TotalSectors    int // total sectors for this range of tracks
}
