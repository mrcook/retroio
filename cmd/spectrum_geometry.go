package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrcook/retroio/spectrum"
	"github.com/mrcook/retroio/spectrum/tap"
	"github.com/mrcook/retroio/spectrum/tzx"
	"github.com/mrcook/retroio/storage"
)

var speccyGeometryCmd = &cobra.Command{
	Use:   "geometry FILE",
	Short: "Read the ZX Spectrum tape geometry",
	Long: `Read the geometry - headers and data tracks/sectors/blocks - from a
ZX Spectrum emulator TZX or TAP file.`,
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
		dskType := mediaType(spectrumMediaType, filename)

		switch dskType {
		case "tap":
			dsk = tap.New(reader)
		case "tzx":
			dsk = tzx.New(reader)
		default:
			fmt.Printf("Unsupported media type: '%s'", dskType)
			return
		}

		if err := dsk.Read(); err != nil {
			fmt.Println("Storage read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		dsk.DisplayGeometry()
	},
}

func init() {
	speccyGeometryCmd.Flags().StringVarP(&spectrumMediaType, "media", "m", "", `Media type, default: file extension`)
	spectrumCmd.AddCommand(speccyGeometryCmd)
}
