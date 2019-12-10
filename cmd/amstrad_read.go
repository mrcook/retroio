package cmd

import (
	"fmt"
	"os"
	"retroio/amstrad/cdt"

	"github.com/spf13/cobra"

	"retroio/spectrum"
	"retroio/storage"
)

var amstradFormat string

var amstradReadCmd = &cobra.Command{
	Use:   "read FILE",
	Short: "Read an Amstrad file",
	Long: `Read all header and data blocks from a CDT file.

NOTE: this storage format is identical to the ZX Spectrum TZX format.`,
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
		dskType := storageType(amstradFormat, filename)

		switch dskType {
		case "cdt":
			dsk = cdt.New(reader)
		default:
			fmt.Printf("Unsupported storage format: '%s'", dskType)
			return
		}

		if err := dsk.Read(); err != nil {
			fmt.Println("Storage read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		dsk.DisplayImageMetadata()
	},
}

func init() {
	amstradReadCmd.Flags().StringVarP(&amstradFormat, "format", "f", "", `Storage format`)
	amstradCmd.AddCommand(amstradReadCmd)
}
