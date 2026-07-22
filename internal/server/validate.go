package server

import (
	"fmt"
	"strings"
)

// ValidateAndRepair normalizes a profile: it trims fields, applies defaults for
// user/port, and resolves conflicting auth methods (password wins over keys, and
// an explicit key wins over an inline key when both are set). It returns an error
// only when a required field is missing or the port is out of range.
func ValidateAndRepair(input Server) (*Server, error) {

	s := input

	s.Name = strings.TrimSpace(s.Name)
	s.Host = strings.TrimSpace(s.Host)
	s.User = strings.TrimSpace(s.User)
	s.Key = strings.TrimSpace(s.Key)
	s.KeyRef = strings.TrimSpace(s.KeyRef)
	s.Password = strings.TrimSpace(s.Password)

	if s.Name == "" {
		return nil, fmt.Errorf("server name is required")
	}

	if s.Host == "" {
		return nil, fmt.Errorf("server host is required")
	}

	if s.User == "" {
		s.User = "root"
	}

	if s.Port == 0 {
		s.Port = 22
	}

	if s.Port < 1 || s.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", s.Port)
	}

	// Password auth takes precedence over any key material.
	if s.Password != "" && (s.Key != "" || s.KeyRef != "") {
		s.Key = ""
		s.KeyRef = ""
	}

	// A stored key reference takes precedence over an inline key path.
	if s.Key != "" && s.KeyRef != "" {
		s.Key = ""
	}

	// Fall back to the conventional default key when no auth method is set.
	if s.Key == "" && s.KeyRef == "" && s.Password == "" {
		s.Key = "~/.ssh/id_rsa"
	}

	return &s, nil
}
