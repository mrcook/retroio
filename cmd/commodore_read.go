package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"retroio/commodore"
	"retroio/commodore/t64"
	"retroio/commodore/tap"
	"retroio/storage"
)

var c64ReadCmd = &cobra.Command{
	Use:                   "read FILE",
	Short:                 "Read a T64 file",
	Long:                  `Read all headers, directories, and data from a T64 image`,
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

		var dsk commodore.Image
		switch path.Ext(filename) {
		case ".tap":
			dsk = tap.New(reader)
		default:
			dsk = t64.New(reader)
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
	commodoreCmd.AddCommand(c64ReadCmd)
}
