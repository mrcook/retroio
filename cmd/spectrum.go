package cmd

import (
	"github.com/spf13/cobra"
)

var (
	spectrumMediaType  string
	spectrumBasListing bool
)

// spectrumCmd represents the spectrum command
var spectrumCmd = &cobra.Command{
	Use:     "spectrum",
	Aliases: []string{"zx"},
	Short:   "System command for the ZX Spectrum",
	Long: `The computer system command for working with disk and tape images for the
Sinclair ZX Spectrum 8-bit home computer.

This is a top-level system command only and requires a sub-command.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(spectrumCmd)
}
