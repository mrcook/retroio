package cmd

import (
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mrcook/retroio/commodore"
)

var commodoreMediaTypeFlag string

// commodoreCmd represents the spectrum command
var commodoreCmd = &cobra.Command{
	Use:     "commodore",
	Aliases: []string{"c64"},
	Short:   "System command for the Commodore C64",
	Long: `The computer system command for working with disk and tape images for the
Sinclair Commodore 8-bit home computers.

This is a top-level system command only and requires a sub-command.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(commodoreCmd)
}

func commodoreDetermineMediaType(filename string) commodore.MediaType {
	if commodoreMediaTypeFlag == "" {
		ext := strings.ToLower(path.Ext(filename))
		commodoreMediaTypeFlag = strings.TrimPrefix(ext, ".")
	}

	switch commodoreMediaTypeFlag {
	case "d64":
		return commodore.D64
	case "d71":
		return commodore.D71
	case "d81":
		return commodore.D81
	case "t64":
		return commodore.T64
	case "tap":
		return commodore.TAP
	default:
		return commodore.Unknown
	}
}
