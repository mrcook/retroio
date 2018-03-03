package cmd

import (
	"fmt"

	"github.com/mrcook/spectrumator/tape"
	"github.com/spf13/cobra"
)

var tzxCmd = &cobra.Command{
	Use:   "tzx FILENAME",
	Short: "Extracts metadata from a TZX file",
	Long:  `Extracts metadata from a given TZX file and prints it to the terminal as formatted text.`,
	Args:  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		tzx := tape.Tzx{}
		if err := tzx.Open(args[0]); err != nil {
			fmt.Println(err)
			return
		}
		defer tzx.Close()

		tzx.Process()
	},
}

func init() {
	rootCmd.AddCommand(tzxCmd)
}
