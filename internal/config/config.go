package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	RepoURL string `yaml:"repo_url"`
	RepoDir string `yaml:"repo_dir"`
}

func GetBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".sshx")
}

func GetRepoDir() string {
	return filepath.Join(GetBaseDir(), "repo")
}

func InitBaseDir() error {
	return os.MkdirAll(GetBaseDir(), 0755)
}
