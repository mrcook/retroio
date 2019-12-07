package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"retroio/spectrum"
	"retroio/spectrum/tap"
	"retroio/spectrum/tzx"
	"retroio/storage"
)

var speccyReadCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a ZX Spectrum file",
	Long:                  `Read all header and data blocks from a TZX or TAP file.`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		normalizeFormatValue()

		filename := args[0]

		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		reader := storage.NewReader(f)

		var dsk spectrum.Image
		dskType := storageType(format, filename)

		switch dskType {
		case "tzx":
			dsk = tzx.New(reader)
		case "tap":
			dsk = tap.New(reader)
		default:
			fmt.Printf("Unsupported tape format: '%s'", dskType)
			return
		}

		if err := dsk.Read(); err != nil {
			fmt.Println("Image read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		dsk.DisplayImageMetadata()
	},
}

func init() {
	speccyReadCmd.Flags().StringVarP(&format, "format", "f", "tzx", `Tape format: TZX or TAP`)
	spectrumCmd.AddCommand(speccyReadCmd)
}
