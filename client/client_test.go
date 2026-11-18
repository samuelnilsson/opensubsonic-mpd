package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser_Success(t *testing.T) {
	// Arrange
	serverInfo := ServerInfo{
		Version:       "1.16.1",
		ServerType:    "OpenSubsonic",
		ServerVersion: "1.0.0",
		OpenSubsonic:  true,
	}
	user := User{
		Username:     "testuser",
		Email:        "test@example.com",
		AdminRole:    true,
		StreamRole:   true,
		DownloadRole: true,
	}

	payload, err := json.Marshal(map[string]any{
		"subsonic-response": map[string]any{
			"status":        "ok",
			"version":       serverInfo.Version,
			"type":          serverInfo.ServerType,
			"serverVersion": serverInfo.ServerVersion,
			"openSubsonic":  serverInfo.OpenSubsonic,
			"user":          user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}))
	defer srv.Close()

	c := New(srv.URL)

	// Act
	result, err := c.GetUser("testuser", "testpass")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Version != "1.16.1" {
		t.Errorf("expected version 1.16.1, got %s", result.Version)
	}
	if result.ServerType != "OpenSubsonic" {
		t.Errorf("expected server type OpenSubsonic, got %s", result.ServerType)
	}
	if result.User.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", result.User.Username)
	}
	if result.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", result.User.Email)
	}
	if !result.User.AdminRole {
		t.Error("expected adminRole true")
	}
}

func TestGetUser_AuthError(t *testing.T) {
	// Arrange
	payload, err := json.Marshal(map[string]any{
		"subsonic-response": map[string]any{
			"status":  "failed",
			"version": "1.16.1",
			"error":   map[string]any{"code": 40, "message": "Wrong username or password"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}))
	defer srv.Close()

	c := New(srv.URL)

	// Act
	_, err = c.GetUser("wrong", "pass")

	// Assert
	if err == nil {
		t.Fatal("expected error for wrong credentials")
	}
}

func TestGetUser_ConnectionError(t *testing.T) {
	// Arrange
	c := New("http://localhost:1") // port unlikely to be open

	// Act
	_, err := c.GetUser("user", "pass")

	// Assert
	if err == nil {
		t.Fatal("expected connection error")
	}
}
