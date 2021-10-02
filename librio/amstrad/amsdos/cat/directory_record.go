package cat

// DirectoryRecord is the displayable data for a directory record.
// This is similar to the CP/M Directory, except each entry merges all record extents.
type DirectoryRecord struct {
	Filename    string
	FileType    string
	RecordCount uint16 // Total record count for all extents of a record

	ReadOnly bool
	Hidden   bool
	Archived bool
}
