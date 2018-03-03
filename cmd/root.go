// Package cmd contains all the CLI commands
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any sub commands
var rootCmd = &cobra.Command{
	Use:     "spectrumator",
	Version: "0.0.0",
	Short:   "Spectrumator: a collection of utility apps for ZX Spectrum hacking",
	Run:     func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
