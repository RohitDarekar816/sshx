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

func LoadServers(repoDir string) ([]Server, error) {

	serverDir := filepath.Join(repoDir, "servers")

	entries, err := os.ReadDir(serverDir)
	if err != nil {
		return nil, err
	}

	var servers []Server

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(serverDir, entry.Name()))
		if err != nil {
			continue
		}

		var s Server
		if err := json.Unmarshal(data, &s); err != nil {
			continue
		}

		servers = append(servers, s)
	}

	return servers, nil
}
