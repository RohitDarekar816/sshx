package server

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadServer(repoDir string, name string) (*Server, error) {

	file := filepath.Join(repoDir, "servers", name+".json")

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var s Server

	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
