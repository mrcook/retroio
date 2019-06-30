package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrcook/tzxit/tzx"
)

var format string

var readCmd = &cobra.Command{
	Use:                   "read TZX/TAP",
	Short:                 "Read a TZX or TAP file",
	Long:                  `Read the metadata from a TZX or TAP file.`,
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

		tape.Read()
		tape.DisplayTapeMetadata()
	},
}

func init() {
	readCmd.Flags().StringVarP(&format, "format", "f", "tzx", `Tape format: TZX or TAP (default: TZX)`)

	rootCmd.AddCommand(readCmd)
}
