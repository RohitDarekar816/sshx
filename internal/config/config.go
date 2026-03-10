package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RepoURL   string `yaml:"repo_url"`
	RepoDir   string `yaml:"repo_dir"`
	UserName  string `yaml:"user_name"`
	UserEmail string `yaml:"user_email"`
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

func ConfigPath() string {
	return filepath.Join(GetBaseDir(), "config.yaml")
}

func LoadConfig() (*Config, error) {

	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {

	if err := InitBaseDir(); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0644)
}
