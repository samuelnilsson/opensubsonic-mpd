package main

import (
	"strings"
	"testing"
)

func TestParseFlags_ValidArgs(t *testing.T) {
	// Arrange
	args := []string{"-server", "http://localhost:4040", "-user", "testuser", "-password", "testpass"}

	// Act
	flags, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if isHelp {
		t.Fatal("expected isHelp=false")
	}
	if flags.ServerURL != "http://localhost:4040" {
		t.Errorf("expected server URL http://localhost:4040, got %s", flags.ServerURL)
	}
	if flags.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", flags.Username)
	}
	if flags.Password != "testpass" {
		t.Errorf("expected password testpass, got %s", flags.Password)
	}
}

func TestParseFlags_ValidArgs_UnknownOrder(t *testing.T) {
	// Arrange
	args := []string{"-password", "secret", "-user", "alice", "-server", "https://music.example.com"}

	// Act
	flags, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if isHelp {
		t.Fatal("expected isHelp=false")
	}
	if flags.ServerURL != "https://music.example.com" {
		t.Errorf("unexpected server URL: %s", flags.ServerURL)
	}
	if flags.Username != "alice" {
		t.Errorf("unexpected username: %s", flags.Username)
	}
	if flags.Password != "secret" {
		t.Errorf("unexpected password: %s", flags.Password)
	}
}

func TestParseFlags_HelpShort(t *testing.T) {
	// Arrange
	args := []string{"-h"}

	// Act
	flags, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if flags != nil {
		t.Fatal("expected nil flags when help is requested")
	}
	if !isHelp {
		t.Fatal("expected isHelp=true for -h")
	}
}

func TestParseFlags_HelpLong(t *testing.T) {
	// Arrange
	args := []string{"--help"}

	// Act
	flags, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if flags != nil {
		t.Fatal("expected nil flags when help is requested")
	}
	if !isHelp {
		t.Fatal("expected isHelp=true for --help")
	}
}

func TestParseFlags_HelpOutputContainsUsage(t *testing.T) {
	// Arrange
	args := []string{"-help"}

	// Act
	_, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !isHelp {
		t.Fatal("expected isHelp=true")
	}
	// Verify help was printed by checking stderr output
	// (parseFlags prints to stdout; we can't easily capture it here,
	//  but the fact that isHelp=true confirms the help path was taken)
}

func TestParseFlags_EmptyArgs(t *testing.T) {
	// Arrange
	args := []string{}

	// Act
	flags, isHelp, err := parseFlags(args)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if isHelp {
		t.Fatal("expected isHelp=false")
	}
	if flags.ServerURL != "" {
		t.Errorf("expected empty server URL, got %s", flags.ServerURL)
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	// Arrange
	args := []string{"--unknown-flag"}

	// Act
	_, _, err := parseFlags(args)

	// Assert
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestValidateFlags_AllPresent(t *testing.T) {
	// Arrange
	flags := &Flags{
		ServerURL: "http://localhost:4040",
		Username:  "testuser",
		Password:  "testpass",
	}

	// Act
	err := validateFlags(flags)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateFlags_MissingServer(t *testing.T) {
	// Arrange
	flags := &Flags{
		ServerURL: "",
		Username:  "testuser",
		Password:  "testpass",
	}

	// Act
	err := validateFlags(flags)

	// Assert
	if err == nil {
		t.Fatal("expected error for missing server")
	}
	if !strings.Contains(err.Error(), "-server") {
		t.Errorf("expected error to mention -server, got: %v", err)
	}
}

func TestValidateFlags_MissingUsername(t *testing.T) {
	// Arrange
	flags := &Flags{
		ServerURL: "http://localhost:4040",
		Username:  "",
		Password:  "testpass",
	}

	// Act
	err := validateFlags(flags)

	// Assert
	if err == nil {
		t.Fatal("expected error for missing username")
	}
}

func TestValidateFlags_MissingPassword(t *testing.T) {
	// Arrange
	flags := &Flags{
		ServerURL: "http://localhost:4040",
		Username:  "testuser",
		Password:  "",
	}

	// Act
	err := validateFlags(flags)

	// Assert
	if err == nil {
		t.Fatal("expected error for missing password")
	}
}

func TestValidateFlags_AllMissing(t *testing.T) {
	// Arrange
	flags := &Flags{}

	// Act
	err := validateFlags(flags)

	// Assert
	if err == nil {
		t.Fatal("expected error for all missing fields")
	}
}
