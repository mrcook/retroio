package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"retroio/tap"
	"retroio/tape"
	"retroio/tzx"
)

var speccyReadCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a TZX or TAP file",
	Long:                  `Read all header and data blocks from a TZX or TAP file.`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		normalizeFormatValue()

		f, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		reader := bufio.NewReader(f)

		var r tape.Tape
		switch format {
		case "tzx":
			r, err = tzx.NewReader(reader)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "tap":
			r = tap.NewReader(reader)
		default:
			fmt.Println("Unsupported tape format.")
			return
		}

		if err := r.ReadBlocks(); err != nil {
			fmt.Println(err)
			return
		}
		r.DisplayTapeMetadata()
	},
}

func init() {
	speccyReadCmd.Flags().StringVarP(&format, "format", "f", "tzx", `Tape format: TZX or TAP`)
	spectrumCmd.AddCommand(speccyReadCmd)
}
