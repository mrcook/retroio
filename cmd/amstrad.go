package cmd

import (
	"github.com/spf13/cobra"
)

// amstradCmd represents the spectrum command
var amstradCmd = &cobra.Command{
	Use:     "amstrad",
	Aliases: []string{"cpc"},
	Short:   "System command for the Amstrad CPC",
	Long: `The computer system command for working with disk and tape images for the
Amstrad CPC 8-bit home computers.

This is a top-level system command only and requires a sub-command.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(amstradCmd)
}
