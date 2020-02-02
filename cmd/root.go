// Copyright (c) 2018-2020 Michael R. Cook. All rights reserved.
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.se.
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any sub commands
var rootCmd = &cobra.Command{
	Use:     "rio",
	Version: "0.10.0",
	Short:   "CLI utility for reading emulator disk and tape images",
	Long: `RetroIO (rio) is a command line utility for reading emulator storage media
(disks and cassette tape images) of home computers from the 1980s.`,
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

func mediaType(media, filename string) string {
	if media == "" {
		media = path.Ext(filename)
	}
	return strings.TrimPrefix(strings.ToLower(media), ".")
}
