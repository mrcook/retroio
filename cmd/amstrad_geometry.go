package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrcook/retroio/amstrad"
	"github.com/mrcook/retroio/amstrad/cdt"
	"github.com/mrcook/retroio/amstrad/dsk"
	"github.com/mrcook/retroio/storage"
)

var amstradMediaType string

var amstradGeometryCmd = &cobra.Command{
	Use:   "geometry FILE",
	Short: "Read the Amstrad disk and tape geometry",
	Long: `Read the geometry - headers and data tracks/sectors/blocks - from an Amstrad
emulator disk or tape file.

NOTE: the CDT geometry is identical to that of the ZX Spectrum TZX format.`,
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

		var disk amstrad.Image
		dskType := mediaType(amstradMediaType, filename)

		switch dskType {
		case "dsk":
			disk = dsk.New(reader)
		case "cdt":
			disk = cdt.New(reader)
		default:
			fmt.Printf("Unsupported media type: '%s'", dskType)
			return
		}

		if err := disk.Read(); err != nil {
			fmt.Println("Media read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		disk.DisplayGeometry()
	},
}

func init() {
	amstradGeometryCmd.Flags().StringVarP(&amstradMediaType, "media", "m", "", `Media type, default: file extension`)
	amstradCmd.AddCommand(amstradGeometryCmd)
}
