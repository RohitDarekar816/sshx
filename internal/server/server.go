package server

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	User     string `json:"user" yaml:"user"`
	Port     int    `json:"port" yaml:"port"`
	Key      string `json:"key,omitempty" yaml:"key,omitempty"`
	KeyRef   string `json:"key_ref,omitempty" yaml:"key_ref,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

func SaveServer(repoDir string, s Server) error {

	serverDir := filepath.Join(repoDir, "servers")

	os.MkdirAll(serverDir, 0755)

	file := filepath.Join(serverDir, s.Name+".yaml")

	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	_ = os.Remove(filepath.Join(serverDir, s.Name+".json"))
	_ = os.Remove(filepath.Join(serverDir, s.Name+".yml"))

	return os.WriteFile(file, data, 0644)
}
