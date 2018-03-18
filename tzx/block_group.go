package tzx

import (
	"fmt"

	"github.com/mrcook/tzxbrowser/tape"
)

// GroupStart
// ID: 21h (33d)
// This block marks the start of a group of blocks which are to be treated as one single
// (composite) block. This is very handy for tapes that use lots of subblocks like Bleepload
// (which may well have over 160 custom loading blocks). You can also give the group a name
// (example 'Bleepload Block 1').
// For each group start block, there must be a group end block. Nesting of groups is not allowed.
type GroupStart struct {
	Length    uint8  // L BYTE  Length of the group name string
	GroupName []byte // CHAR[L] Group name in ASCII format (please keep it under 30 characters long)
}

func (g *GroupStart) Process(file *tape.File) {
	g.Length, _ = file.ReadByte()

	for _, b := range file.ReadBytes(int(g.Length)) {
		g.GroupName = append(g.GroupName, b)
	}
}

func (g GroupStart) Id() int {
	return 33
}

func (g GroupStart) Name() string {
	return "Group Start"
}

// Metadata returns a human readable string of the block data
func (g GroupStart) Metadata() string {
	return fmt.Sprintf("> %-19s : %s", g.Name(), g.GroupName)
}

// GroupEnd
// ID: 22h (34d)
// This indicates the end of a group. This block has no body.
type GroupEnd struct{}

func (g *GroupEnd) Process(file *tape.File) {}

func (g GroupEnd) Id() int {
	return 34
}

func (g GroupEnd) Name() string {
	return "Group End"
}

// Metadata returns a human readable string of the block data
func (g GroupEnd) Metadata() string {
	return fmt.Sprintf("> %s", g.Name())
}
