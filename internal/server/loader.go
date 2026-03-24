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
		server, ext := processEntry(serverDir, entry, byName, sourceExt)
		if server == nil {
			continue
		}

		if _, ok := byName[server.Name]; !ok {
			order = append(order, server.Name)
		}

		byName[server.Name] = *server
		sourceExt[server.Name] = ext
	}

	return buildServerList(byName, order), nil
}

func processEntry(serverDir string, entry os.DirEntry, byName map[string]Server, sourceExt map[string]string) (*Server, string) {
	if entry.IsDir() {
		return nil, ""
	}

	ext := filepath.Ext(entry.Name())
	if !isValidServerFile(ext) {
		return nil, ""
	}

	server, err := loadAndValidateServer(serverDir, entry.Name(), ext)
	if err != nil {
		return nil, ""
	}

	if shouldSkipForDuplicate(server.Name, ext, sourceExt) {
		return nil, ""
	}

	return server, ext
}

func isValidServerFile(ext string) bool {
	return ext == ".json" || ext == ".yaml" || ext == ".yml"
}

func loadAndValidateServer(serverDir, filename, ext string) (*Server, error) {
	fullPath := filepath.Join(serverDir, filename)
	s, err := loadServerFile(fullPath)
	if err != nil {
		return nil, err
	}

	repaired, _ := ValidateAndRepair(*s)
	return repaired, nil
}

func shouldSkipForDuplicate(name, ext string, sourceExt map[string]string) bool {
	existingExt, ok := sourceExt[name]
	if !ok {
		return false
	}

	isYamlExisting := existingExt == ".yaml" || existingExt == ".yml"
	isJsonNew := ext == ".json"
	return isYamlExisting && isJsonNew
}

func buildServerList(byName map[string]Server, order []string) []Server {
	servers := make([]Server, 0, len(byName))
	for _, name := range order {
		if s, ok := byName[name]; ok {
			servers = append(servers, s)
		}
	}
	return servers
}

func RemoveServer(repoDir, name string) error {

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
