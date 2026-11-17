package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("h", false, "show help and exit")
	longHelp := flag.Bool("help", false, "show help and exit")
	flag.Parse()

	if *help || *longHelp {
		fmt.Fprintf(os.Stderr, `Usage: opensubsonic-mpd [options]

Options:
  -h, --help      show this help message

`)
		os.Exit(0)
	}

	fmt.Println("hello world")
}
