package main

import (
	"bytes"
	"fmt"
	"testing"

	"opensubsonic-mpd/client"
)

// fakeRunner implements the Runner interface for testing Run().
type fakeRunner struct {
	result *client.AuthResult
	err    error
}

func (f *fakeRunner) GetUser(username, password string) (*client.AuthResult, error) {
	return f.result, f.err
}

func (f *fakeRunner) GetMusicDirectory(id string) (*client.DirectoryResult, error) {
	return &client.DirectoryResult{
		Directory: &client.Directory{
			Name: "Music",
			Child: []client.Child{
				{IsDir: true, Title: "Albums"},
				{IsDir: true, Title: "Artists"},
			},
		},
	}, nil
}

func TestRun_Success(t *testing.T) {
	// Arrange
	args := RunArgs{
		ServerURL: "http://localhost:4040",
		Username:  "testuser",
		Password:  "testpass",
	}

	fake := &fakeRunner{
		result: &client.AuthResult{
			ServerInfo: client.ServerInfo{
				Version:       "1.16.1",
				ServerType:    "OpenSubsonic",
				ServerVersion: "1.0.0",
				OpenSubsonic:  true,
			},
			User: &client.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
	}

	var buf bytes.Buffer

	// Act
	err := runWith(&buf, args, fake)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := `Authenticated as testuser
Server: OpenSubsonic 1.0.0
API version: 1.16.1
OpenSubsonic server

Root directory (Music):
  📁 Albums
  📁 Artists
`
	if buf.String() != expected {
		t.Errorf("unexpected output:\nexpected:\n%s\n\nactual:\n%s", expected, buf.String())
	}
}

func TestRun_Success_NoOpenSubsonic(t *testing.T) {
	// Arrange
	args := RunArgs{
		ServerURL: "http://localhost:4040",
		Username:  "testuser",
		Password:  "testpass",
	}

	fake := &fakeRunner{
		result: &client.AuthResult{
			ServerInfo: client.ServerInfo{
				Version:       "1.16.1",
				ServerType:    "Navidrome",
				ServerVersion: "0.50.0",
				OpenSubsonic:  false,
			},
			User: &client.User{
				Username: "testuser",
			},
		},
	}

	var buf bytes.Buffer

	// Act
	err := runWith(&buf, args, fake)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := `Authenticated as testuser
Server: Navidrome 0.50.0
API version: 1.16.1

Root directory (Music):
  📁 Albums
  📁 Artists
`
	if buf.String() != expected {
		t.Errorf("unexpected output:\nexpected:\n%s\n\nactual:\n%s", expected, buf.String())
	}
}

func TestRun_AuthError(t *testing.T) {
	// Arrange
	args := RunArgs{
		ServerURL: "http://localhost:4040",
		Username:  "wronguser",
		Password:  "wrongpass",
	}

	fake := &fakeRunner{
		err: fmt.Errorf("server error 40: Wrong username or password"),
	}

	var buf bytes.Buffer

	// Act
	err := runWith(&buf, args, fake)

	// Assert
	if err == nil {
		t.Fatal("expected error for auth failure")
	}

	// The error should contain a descriptive message
	expectedPrefix := "authentication failed:"
	if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected error prefix %q, got: %v", expectedPrefix, err)
	}

	// Nothing should be written on error
	if buf.Len() > 0 {
		t.Errorf("expected no output on error, got: %q", buf.String())
	}
}
