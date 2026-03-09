package server

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Key      string `json:"key,omitempty"`
	Password string `json:"password,omitempty"`
}

func SaveServer(repoDir string, s Server) error {

	serverDir := filepath.Join(repoDir, "servers")

	os.MkdirAll(serverDir, 0755)

	file := filepath.Join(serverDir, s.Name+".json")

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
