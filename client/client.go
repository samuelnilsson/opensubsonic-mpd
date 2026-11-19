package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiVersion = "1.16.1"
	clientName = "opensubsonic-mpd"
)

// ServerInfo holds metadata returned by the getUser endpoint.
type ServerInfo struct {
	Version       string `json:"version"`
	ServerType    string `json:"type,omitempty"`
	ServerVersion string `json:"serverVersion,omitempty"`
	OpenSubsonic  bool   `json:"openSubsonic,omitempty"`
}

// Server wraps a base URL and provides OpenSubsonic API methods.
type Server struct {
	BaseURL string
	HTTP    *http.Client
}

// New creates a new Server with the given base URL (e.g. "http://localhost:4040").
func New(baseURL string) *Server {
	return &Server{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{},
	}
}

// GetUser calls the /rest/getUser endpoint and verifies the credentials
// by returning the authenticated user's info.
func (s *Server) GetUser(username, password string) (*AuthResult, error) {
	endpoint := fmt.Sprintf("%s/rest/getUser.view", s.BaseURL)

	params := url.Values{}
	params.Set("v", apiVersion)
	params.Set("c", clientName)
	params.Set("f", "json")
	params.Set("u", username)
	params.Set("p", password)
	params.Set("username", username)

	resp, err := s.HTTP.Get(endpoint + "?" + params.Encode())
	if err != nil {
		msg := err.Error()
		msg = strings.ReplaceAll(msg, "p="+params.Get("p"), "p=***")
		return nil, fmt.Errorf("connection failed: %s", msg)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		SubsonicResponse struct {
			Status        string    `json:"status"`
			Version       string    `json:"version"`
			ServerType    string    `json:"type,omitempty"`
			ServerVersion string    `json:"serverVersion,omitempty"`
			OpenSubsonic  bool      `json:"openSubsonic,omitempty"`
			User          *User     `json:"user,omitempty"`
			Error         *APIError `json:"error,omitempty"`
		} `json:"subsonic-response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.SubsonicResponse.Error != nil {
		e := result.SubsonicResponse.Error
		return nil, fmt.Errorf("server error %d: %s", e.Code, e.Message)
	}

	if result.SubsonicResponse.Status != "ok" {
		return nil, fmt.Errorf("unexpected status: %s", result.SubsonicResponse.Status)
	}

	return &AuthResult{
		ServerInfo: ServerInfo{
			Version:       result.SubsonicResponse.Version,
			ServerType:    result.SubsonicResponse.ServerType,
			ServerVersion: result.SubsonicResponse.ServerVersion,
			OpenSubsonic:  result.SubsonicResponse.OpenSubsonic,
		},
		User: result.SubsonicResponse.User,
	}, nil
}

// AuthResult holds the server and user info returned after successful auth.
type AuthResult struct {
	ServerInfo
	User *User
}

// User represents an authenticated OpenSubsonic user.
type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	AdminRole    bool   `json:"adminRole"`
	StreamRole   bool   `json:"streamRole"`
	DownloadRole bool   `json:"downloadRole"`
	PlaylistRole bool   `json:"playlistRole"`
	ShareRole    bool   `json:"shareRole"`
}

// GetMusicDirectory calls the /rest/getMusicDirectory endpoint and returns
// the contents of the specified music directory. Pass an empty string for "id"
// to get the root directory.
func (s *Server) GetMusicDirectory(id string) (*DirectoryResult, error) {
	endpoint := fmt.Sprintf("%s/rest/getMusicDirectory.view", s.BaseURL)

	params := url.Values{}
	params.Set("v", apiVersion)
	params.Set("c", clientName)
	params.Set("f", "json")
	params.Set("id", id)

	resp, err := s.HTTP.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		SubsonicResponse struct {
			Status   string `json:"status"`
			Version  string `json:"version"`
			Error    *APIError `json:"error,omitempty"`
			Directory *Directory `json:"directory,omitempty"`
		} `json:"subsonic-response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.SubsonicResponse.Error != nil {
		e := result.SubsonicResponse.Error
		return nil, fmt.Errorf("server error %d: %s", e.Code, e.Message)
	}

	if result.SubsonicResponse.Status != "ok" {
		return nil, fmt.Errorf("unexpected status: %s", result.SubsonicResponse.Status)
	}

	if result.SubsonicResponse.Directory == nil {
		return nil, fmt.Errorf("directory not found")
	}

	return &DirectoryResult{Directory: result.SubsonicResponse.Directory}, nil
}

// APIError represents an error response from the OpenSubsonic API.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Directory represents a music directory returned by the server.
type Directory struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Child   []Child `json:"child"`
}

// Child represents a single item within a directory (file or subdirectory).
type Child struct {
	ID              string  `json:"id"`
	Parent          string  `json:"parent,omitempty"`
	IsDir           bool    `json:"isDir"`
	Title           string  `json:"title"`
	Album           string  `json:"album,omitempty"`
	Artist          string  `json:"artist,omitempty"`
	Track           int     `json:"track,omitempty"`
	Year            int     `json:"year,omitempty"`
	Genre           string  `json:"genre,omitempty"`
	CoverArt        string  `json:"coverArt,omitempty"`
	Size            int64   `json:"size,omitempty"`
	ContentType     string  `json:"contentType,omitempty"`
	Suffix          string  `json:"suffix,omitempty"`
	Duration        int     `json:"duration,omitempty"`
	BitRate         int     `json:"bitRate,omitempty"`
	Path            string  `json:"path,omitempty"`
	AlbumID         string  `json:"albumId,omitempty"`
	ArtistID        string  `json:"artistId,omitempty"`
	Type            string  `json:"type,omitempty"`
	UserRating      int     `json:"userRating,omitempty"`
	PlayCount       int     `json:"playCount,omitempty"`
	DiscNumber      int     `json:"discNumber,omitempty"`
	Created         string  `json:"created,omitempty"`
	Starred         string  `json:"starred,omitempty"`
	Played          string  `json:"played,omitempty"`
	BookmarkPos     int     `json:"bookmarkPosition,omitempty"`
	OriginalWidth   int     `json:"originalWidth,omitempty"`
	OriginalHeight  int     `json:"originalHeight,omitempty"`
	BPM             int     `json:"bpm,omitempty"`
	MusicBrainzID   string  `json:"musicBrainzId,omitempty"`
	ExplicitStatus  string  `json:"explicitStatus,omitempty"`
}

// DirectoryResult wraps the directory returned by GetMusicDirectory.
type DirectoryResult struct {
	*Directory
}
