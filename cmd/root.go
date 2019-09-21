// Package cmd contains all the CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any sub commands
var rootCmd = &cobra.Command{
	Use:     "tzxit",
	Version: "0.5.0",
	Short:   "tzxit: a CLI based ZX Spectrum cassette tape utility",
	Long: `tzxit is a command line utility for reading and writing ZX Spectrum
tape cassette images using the TZX specification.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(cmd.ValidArgs) == 0 {
			_ = cmd.Help()
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
