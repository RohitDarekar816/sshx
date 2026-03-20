package totp

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SecretFile struct {
	Version int    `yaml:"version"`
	Cipher  string `yaml:"cipher"`
	KDF     string `yaml:"kdf"`
	Data    string `yaml:"data"`
}

type Config struct {
	RepoDir string
	Email   string
}

func SecretPath(cfg Config) string {
	dir := filepath.Join(cfg.RepoDir, "totp")
	return filepath.Join(dir, sanitize(cfg.Email)+".yaml")
}

func HasSecret(cfg Config) bool {
	_, err := os.Stat(SecretPath(cfg))
	return err == nil
}

func SaveEncryptedSecret(cfg Config, encrypted string) error {
	dir := filepath.Join(cfg.RepoDir, "totp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(SecretFile{
		Version: 1,
		Cipher:  "AES-256-GCM",
		KDF:     "scrypt",
		Data:    encrypted,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(SecretPath(cfg), data, 0644)
}

func LoadEncryptedSecret(cfg Config) (string, error) {
	path := SecretPath(cfg)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var f SecretFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return "", err
	}

	return f.Data, nil
}

func sanitize(value string) string {

	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}

	return b.String()
}
