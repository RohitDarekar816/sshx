package server

import (
	"fmt"
	"strings"
)

func ValidateAndRepair(input Server) (*Server, error) {

	s := input
	changed := false

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
		changed = true
	}

	if s.Port == 0 {
		s.Port = 22
		changed = true
	}

	if s.Port < 1 || s.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", s.Port)
	}

	if s.Password != "" && (s.Key != "" || s.KeyRef != "") {
		s.Key = ""
		s.KeyRef = ""
		changed = true
	}

	if s.Key != "" && s.KeyRef != "" {
		s.Key = ""
		changed = true
	}

	if s.Key == "" && s.KeyRef == "" && s.Password == "" {
		s.Key = "~/.ssh/id_rsa"
		changed = true
	}

	if changed {
		return &s, nil
	}

	return &s, nil
}
