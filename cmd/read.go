package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mrcook/tzxit/tape"
	"github.com/mrcook/tzxit/tzx"
)

var format string

var readCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a TZX or TAP file",
	Long:                  `Read the metadata from a TZX or TAP file.`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		var r tape.Tape

		normalizeFormatValue()
		if format == "tap" {
			fmt.Println("TAP reading not yet supported")
			return
		} else if format == "tzx" {
			if r, err = tzx.NewReader(file); err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			fmt.Println("Unsupported tape format")
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
	readCmd.Flags().StringVarP(&format, "format", "f", "tzx", `Tape format: TZX or TAP`)
	rootCmd.AddCommand(readCmd)
}

func normalizeFormatValue() {
	format = strings.ToLower(format)
}
