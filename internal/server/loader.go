package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadServer(repoDir, name string) (*Server, error) {

	serverDir := filepath.Join(repoDir, "servers")
	candidates := []string{
		filepath.Join(serverDir, name+".yaml"),
		filepath.Join(serverDir, name+".yml"),
		filepath.Join(serverDir, name+".json"),
	}

	for _, file := range candidates {
		s, err := loadServerFile(file)
		if err == nil {
			repaired, _ := ValidateAndRepair(*s)
			return repaired, nil
		}
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return nil, os.ErrNotExist
}

func LoadServers(repoDir string) ([]Server, error) {

	serverDir := filepath.Join(repoDir, "servers")

	entries, err := os.ReadDir(serverDir)
	if err != nil {
		return nil, err
	}

	byName := map[string]Server{}
	sourceExt := map[string]string{}
	order := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		fullPath := filepath.Join(serverDir, entry.Name())
		s, err := loadServerFile(fullPath)
		if err != nil {
			continue
		}

		repaired, _ := ValidateAndRepair(*s)
		s = repaired

		if _, ok := byName[s.Name]; !ok {
			order = append(order, s.Name)
		}

		// Prefer YAML over JSON when duplicates exist.
		if existingExt, ok := sourceExt[s.Name]; ok {
			if (existingExt == ".yaml" || existingExt == ".yml") && ext == ".json" {
				continue
			}
		}

		byName[s.Name] = *s
		sourceExt[s.Name] = ext
	}

	servers := make([]Server, 0, len(byName))
	for _, name := range order {
		if s, ok := byName[name]; ok {
			servers = append(servers, s)
		}
	}

	return servers, nil
}

func RemoveServer(repoDir string, name string) error {

	serverDir := filepath.Join(repoDir, "servers")
	candidates := []string{
		filepath.Join(serverDir, name+".yaml"),
		filepath.Join(serverDir, name+".yml"),
		filepath.Join(serverDir, name+".json"),
	}

	var removed bool
	var lastErr error
	for _, file := range candidates {
		if err := os.Remove(file); err == nil {
			removed = true
		} else if !os.IsNotExist(err) {
			lastErr = err
		}
	}

	if removed {
		return nil
	}

	if lastErr != nil {
		return lastErr
	}

	return os.ErrNotExist
}

func loadServerFile(path string) (*Server, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s Server
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &s); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}
	}

	if s.Name == "" {
		base := filepath.Base(path)
		s.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return &s, nil
}
