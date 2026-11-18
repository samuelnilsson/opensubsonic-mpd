package main

import (
	"bytes"
	"flag"
	"fmt"
)

// Flags holds the parsed command-line arguments.
type Flags struct {
	ServerURL string
	Username  string
	Password  string
}

// parseFlags parses the command-line flags.
// Returns the parsed Flags, a boolean indicating if help was requested,
// and an error if the flags could not be parsed.
func parseFlags(args []string) (*Flags, bool, error) {
	fs := flag.NewFlagSet("opensubsonic-mpd", flag.ContinueOnError)
	serverURL := fs.String("server", "", "OpenSubsonic server URL (e.g. http://localhost:4040)")
	username := fs.String("user", "", "Username for authentication")
	password := fs.String("password", "", "Password for authentication")
	showHelp := fs.Bool("h", false, "show help and exit")
	longHelp := fs.Bool("help", false, "show help and exit")

	var buf bytes.Buffer
	fs.SetOutput(&buf)

	if err := fs.Parse(args); err != nil {
		return nil, false, err
	}

	if *showHelp || *longHelp {
		fmt.Fprintf(&buf, `Usage: opensubsonic-mpd [options]

Connect to an OpenSubsonic server and expose it as an MPD-compatible audio player.

Options:
  -server URL     OpenSubsonic server URL (required, e.g. http://localhost:4040)
  -user NAME      Username for authentication (required)
  -password PASS  Password for authentication (required)
  -h, --help      show this help message

`)
		fmt.Print(buf.String())
		return nil, true, nil
	}

	return &Flags{
		ServerURL: *serverURL,
		Username:  *username,
		Password:  *password,
	}, false, nil
}

// validateFlags checks that all required arguments are provided.
func validateFlags(f *Flags) error {
	if f.ServerURL == "" || f.Username == "" || f.Password == "" {
		return fmt.Errorf("-server, -user, and -password are all required")
	}
	return nil
}
