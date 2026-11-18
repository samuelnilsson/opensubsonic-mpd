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
}

// Run executes the core application logic and returns an error on failure.
// It writes status information to w (use os.Stdout for production).
func Run(w io.Writer, args RunArgs) error {
	srv := client.New(args.ServerURL)
	result, err := srv.GetUser(args.Username, args.Password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Fprintf(w, "Authenticated as %s\n", result.User.Username)
	fmt.Fprintf(w, "Server: %s %s\n", result.ServerType, result.ServerVersion)
	fmt.Fprintf(w, "API version: %s\n", result.Version)
	if result.OpenSubsonic {
		fmt.Fprintln(w, "OpenSubsonic server")
	}

	return nil
}

// RunStdout is the production entry point — delegates to Run with os.Stdout.
func RunStdout(args RunArgs) error {
	return Run(os.Stdout, args)
}

// runWith is the internal implementation that accepts an injectable Runner.
// RunStdout and tests both delegate to this function.
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

	return nil
}
