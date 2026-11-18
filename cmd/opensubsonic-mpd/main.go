package main

import (
	"fmt"
	"os"
)

func main() {
	flags, isHelp, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isHelp {
		return // parseFlags already printed help and called exit(0)
	}

	if err := validateFlags(flags); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		fmt.Fprintln(os.Stderr, "Run with -help for usage information")
		os.Exit(1)
	}

	if err := RunStdout(RunArgs{
		ServerURL: flags.ServerURL,
		Username:  flags.Username,
		Password:  flags.Password,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
