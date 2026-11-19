package main

import (
	"fmt"
	"io"
	"os"

	"opensubsonic-mpd/client"
)

// RunArgs holds the validated command-line arguments passed to run().
type RunArgs struct {
	ServerURL string
	Username  string
	Password  string
}

// Runner abstracts the client dependency so it can be mocked in tests.
type Runner interface {
	GetUser(username, password string) (*client.AuthResult, error)
	GetMusicFolders(username, password string) (*client.MusicFoldersResult, error)
	GetMusicDirectory(id, username, password string) (*client.DirectoryResult, error)
}

// Run executes the core application logic and returns an error on failure.
// It writes status information to w (use os.Stdout for production).
func Run(w io.Writer, args RunArgs) error {
	srv := client.New(args.ServerURL)
	return runWith(w, args, srv)
}

// RunStdout is the production entry point — delegates to Run with os.Stdout.
func RunStdout(args RunArgs) error {
	return Run(os.Stdout, args)
}

// runWith is the internal implementation that accepts an injectable Runner.
// Run and tests both delegate to this function.
func runWith(w io.Writer, args RunArgs, runner Runner) error {
	result, err := runner.GetUser(args.Username, args.Password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Fprintf(w, "Authenticated as %s\n", result.User.Username)
	fmt.Fprintf(w, "Server: %s %s\n", result.ServerType, result.ServerVersion)
	fmt.Fprintf(w, "API version: %s\n", result.Version)
	if result.OpenSubsonic {
		fmt.Fprintln(w, "OpenSubsonic server")
	}

	// List root directory
	folders, err := runner.GetMusicFolders(args.Username, args.Password)
	if err != nil {
		fmt.Fprintf(w, "Warning: failed to get music folders: %v\n", err)
	} else if len(folders.Folders) == 0 {
		fmt.Fprintln(w, "Warning: no music folders configured")
	} else {
		folder := folders.Folders[0]
		fmt.Fprintf(w, "\nMusic folders:\n")
		for _, f := range folders.Folders {
			fmt.Fprintf(w, "  [id=%d] %s\n", f.ID, f.Name)
		}

		rootDir, err := runner.GetMusicDirectory(fmt.Sprintf("%d", folder.ID), args.Username, args.Password)
		if err != nil {
			fmt.Fprintf(w, "Warning: failed to get music folder contents: %v\n", err)
		} else {
			fmt.Fprintf(w, "\nContents of %q:\n", rootDir.Name)
			for _, item := range rootDir.Child {
				indent := ""
				if item.IsDir {
					indent = "📁 "
				} else {
					indent = "🎵 "
				}
				fmt.Fprintf(w, "  %s%s\n", indent, item.Title)
			}
		}
	}

	return nil
}
