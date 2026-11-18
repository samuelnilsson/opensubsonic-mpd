package main

import (
	"flag"
	"fmt"
	"os"

	"opensubsonic-mpd/client"
)

func main() {
	serverURL := flag.String("server", "", "OpenSubsonic server URL (e.g. http://localhost:4040)")
	username := flag.String("user", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")
	showHelp := flag.Bool("h", false, "show help and exit")
	longHelp := flag.Bool("help", false, "show help and exit")
	flag.Parse()

	if *showHelp || *longHelp {
		fmt.Fprintf(os.Stderr, `Usage: opensubsonic-mpd [options]

Connect to an OpenSubsonic server and expose it as an MPD-compatible audio player.

Options:
  -server URL     OpenSubsonic server URL (required, e.g. http://localhost:4040)
  -user NAME      Username for authentication (required)
  -password PASS  Password for authentication (required)
  -h, --help      show this help message

`)
		os.Exit(0)
	}

	if *serverURL == "" || *username == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Error: -server, -user, and -password are all required")
		fmt.Fprintln(os.Stderr, "Run with -help for usage information")
		os.Exit(1)
	}

	srv := client.New(*serverURL)
	result, err := srv.GetUser(*username, *password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Authentication failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Authenticated as %s\n", result.User.Username)
	fmt.Printf("Server: %s %s\n", result.ServerType, result.ServerVersion)
	fmt.Printf("API version: %s\n", result.Version)
	if result.OpenSubsonic {
		fmt.Println("OpenSubsonic server")
	}
}
