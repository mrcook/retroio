package cmd

import (
	"fmt"
	"os"

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
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		r, err := tzx.NewReader(file)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := r.ReadBlocks(); err != nil {
			fmt.Println(err.Error())
			return
		}
		r.DisplayTapeMetadata()
	},
}

func init() {
	readCmd.Flags().StringVarP(&format, "format", "f", "tzx", `Tape format: TZX or TAP (default: TZX)`)

	rootCmd.AddCommand(readCmd)
}
