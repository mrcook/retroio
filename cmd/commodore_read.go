package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"retroio/commodore"
	"retroio/commodore/t64"
	"retroio/commodore/tap"
	"retroio/storage"
)

var commodoreFormat string

var c64ReadCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a T64 file",
	Long:                  `Read all headers, directories, and data from a T64 image`,
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

		var dsk commodore.Image
		dskType := storageType(commodoreFormat, filename)

		switch dskType {
		case "t64":
			dsk = t64.New(reader)
		case "tap":
			dsk = tap.New(reader)
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
	c64ReadCmd.Flags().StringVarP(&commodoreFormat, "format", "f", "", `Storage format`)
	commodoreCmd.AddCommand(c64ReadCmd)
}
