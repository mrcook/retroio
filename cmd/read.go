package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrcook/tzxit/tzx"
)

var readCmd = &cobra.Command{
	Use:                   "read TZX",
	Short:                 "Read a TZX file",
	Long:                  `Read the metadata from a TZX file.`,
	Aliases:               []string{"r"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		tape := tzx.Tzx{}
		if err := tape.Open(args[0]); err != nil {
			fmt.Println(err)
			return
		}
		defer tape.Close()

		tape.Process()
		tape.DisplayTapeMetadata()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
