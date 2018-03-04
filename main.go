package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrcook/tzxbrowser/tape"
)

const version = "0.0.1"

func init() {
	showVersion := flag.Bool("v", false, "Prints the current version")

	flag.Usage = func() {
		fmt.Println("TZX Browser - ZX Spectrum tape file parser")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("  %s [OPTIONS] /path/to/tape.tzx", os.Args[0])
		fmt.Println()
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("  -h, --help       Show this message")
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("TZX Browser v%s\n", version)
		os.Exit(0)
	}

	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Println("ERROR, you must specify a filename.")
		fmt.Println()
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	tzx := tape.Tzx{}
	if err := tzx.Open(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}
	defer tzx.Close()

	tzx.Process()
	tzx.DisplayTapeMetadata()
}
