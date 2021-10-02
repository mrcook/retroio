package disk

// Directory File Type Labels
var fileTypes = []FileType{
	{0, "DEL", false, false, true, "Deleted"},
	{1, "SEQ", false, false, true, "Sequential"},
	{2, "PRG", false, false, true, "Program"},
	{3, "USR", false, false, true, "User"},
	{4, "REL", false, false, true, "Relative"},

	// D81 ONLY
	{5, "CBM", false, false, true, "Partition/Sub-directory"},
}
