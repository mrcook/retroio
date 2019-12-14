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

var (
	spectrumFormat     string
	spectrumBasListing bool
)

var speccyReadCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a ZX Spectrum file",
	Long:                  `Read all header and data blocks from a TZX or TAP file.`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		reader := storage.NewReader(f)

		var dsk spectrum.Image
		dskType := storageType(spectrumFormat, filename)

		switch dskType {
		case "tap":
			dsk = tap.New(reader)
		case "tzx":
			dsk = tzx.New(reader)
		default:
			fmt.Printf("Unsupported storage format: '%s'", dskType)
			return
		}

		if err := dsk.Read(); err != nil {
			fmt.Println("Storage read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		if spectrumBasListing {
			dsk.ListBasicPrograms()
		} else {
			dsk.DisplayImageMetadata()
		}
	},
}

func init() {
	speccyReadCmd.Flags().StringVarP(&spectrumFormat, "format", "f", "", `Storage format`)
	speccyReadCmd.Flags().BoolVar(&spectrumBasListing, "bas", false, `BASIC program listing`)
	spectrumCmd.AddCommand(speccyReadCmd)
}
