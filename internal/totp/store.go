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

// SecretPath returns the path to the TOTP secret file for the given repository directory and email.
// Issue #10: SonarQube - Group together these consecutive parameters of the same type.
func SecretPath(repoDir, email string) string {

	dir := filepath.Join(repoDir, "totp")
	return filepath.Join(dir, sanitize(email)+".yaml")
}

func HasSecret(repoDir, email string) bool {
	_, err := os.Stat(SecretPath(repoDir, email))
	return err == nil
}

func SaveEncryptedSecret(repoDir, email, encrypted string) error {

	dir := filepath.Join(repoDir, "totp")
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

	return os.WriteFile(SecretPath(repoDir, email), data, 0644)
}

func LoadEncryptedSecret(repoDir, email string) (string, error) {

	path := SecretPath(repoDir, email)
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
