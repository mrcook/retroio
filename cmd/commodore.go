package cmd

import (
	"github.com/spf13/cobra"
)

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
