package trd

// Directory (2048 bytes)
//
//     Offset | Type          | Length | Description
//         0    DIRENTRY #0       16     header of 1st file
//        16    DIRENTRY #1       16     header of 2nd file
//        32    DIRENTRY #2       16     header of 3rd file
//      ....
//      2032    DIRENTRY #127     16     header of 128th file
//
// The directory always contains 128 entries, but the list ends when the first
// byte of the directory entry (= of the file name) is #0. Therefore, any
// following entries are invalid and may contain random data.
type Directory struct {
	Entries [128]DirEntry
}
