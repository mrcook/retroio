package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"retroio/amstrad"
	"retroio/amstrad/cdt"
	"retroio/amstrad/dsk"
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

		var disk amstrad.Image
		dskType := storageType(amstradFormat, filename)

		switch dskType {
		case "dsk":
			disk = dsk.New(reader)
		case "cdt":
			disk = cdt.New(reader)
		default:
			fmt.Printf("Unsupported storage format: '%s'", dskType)
			return
		}

		if err := disk.Read(); err != nil {
			fmt.Println("Storage read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		disk.DisplayImageMetadata()
	},
}

func init() {
	amstradReadCmd.Flags().StringVarP(&amstradFormat, "format", "f", "", `Storage format`)
	amstradCmd.AddCommand(amstradReadCmd)
}
