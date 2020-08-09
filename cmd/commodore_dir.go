package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrcook/retroio/commodore"
	"github.com/mrcook/retroio/commodore/d64"
	"github.com/mrcook/retroio/commodore/d71"
	"github.com/mrcook/retroio/commodore/d81"
	"github.com/mrcook/retroio/storage"
)

var commodoreCommandDir = &cobra.Command{
	Use:                   "dir FILE",
	Short:                 "Displays the directory of a Commodore disk image",
	Long:                  `Performs a directory listing for Commodore D64, D71, and D81 disk image files.`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		reader, err := storage.NewReaderFromFile(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		diskSize := uint32(reader.FileSize)

		mediaType := commodoreDetermineMediaType(reader.Filename)
		if mediaType == commodore.Unknown {
			fmt.Printf("unknown media type for %s", reader.Filename)
			return
		}

		var dsk commodore.Image

		switch mediaType {
		case commodore.D64:
			if dsk, err = d64.New(reader, diskSize); err != nil {
				fmt.Println(err)
				return
			}
		case commodore.D71:
			if dsk, err = d71.New(reader, diskSize); err != nil {
				fmt.Println(err)
				return
			}
		case commodore.D81:
			if dsk, err = d81.New(reader, diskSize); err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Print("unsupported media type for this command")
			return
		}

		if err := dsk.Read(); err != nil {
			fmt.Println("Media read error!")
			fmt.Println(err)
			os.Exit(1)
		}

		dsk.CommandDir()
	},
}

func init() {
	commodoreCommandDir.Flags().StringVarP(&commodoreMediaTypeFlag, "media", "m", "", `Media type, default: file extension`)
	commodoreCmd.AddCommand(commodoreCommandDir)
}
