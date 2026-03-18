package keys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/crypto"
	"gopkg.in/yaml.v3"
)

type KeyFile struct {
	Version int    `yaml:"version"`
	Cipher  string `yaml:"cipher"`
	KDF     string `yaml:"kdf"`
	Data    string `yaml:"data"`
}

func KeyPath(repoDir, name string) string {
	return filepath.Join(repoDir, "keys", name+".yaml")
}

func EncryptKeyFile(repoDir, name string, filePath, passphrase string) error {

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(string(raw), passphrase)
	if err != nil {
		return err
	}

	return SaveEncryptedKey(repoDir, name, encrypted)
}

func SaveEncryptedKey(repoDir, name, encrypted string) error {

	dir := filepath.Join(repoDir, "keys")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(KeyFile{
		Version: 1,
		Cipher:  "AES-256-GCM",
		KDF:     "scrypt",
		Data:    encrypted,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(KeyPath(repoDir, name), data, 0644)
}

func LoadEncryptedKey(repoDir, name string) (string, error) {

	data, err := os.ReadFile(KeyPath(repoDir, name))
	if err != nil {
		return "", err
	}

	var f KeyFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return "", err
	}

	if f.Data == "" {
		return "", fmt.Errorf("empty encrypted key")
	}

	return f.Data, nil
}

func DecryptKey(repoDir, name, passphrase string) ([]byte, error) {

	encrypted, err := LoadEncryptedKey(repoDir, name)
	if err != nil {
		return nil, err
	}

	plain, err := crypto.Decrypt(encrypted, passphrase)
	if err != nil {
		return nil, err
	}

	return []byte(plain), nil
}

func ListKeys(repoDir string) ([]string, error) {

	dir := filepath.Join(repoDir, "keys")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		if name != "" {
			names = append(names, name)
		}
	}

	return names, nil
}

func RemoveKey(repoDir, name string) error {
	return os.Remove(KeyPath(repoDir, name))
}

func WriteTempKey(data []byte) (string, error) {

	f, err := os.CreateTemp("", "sshx-key-*")
	if err != nil {
		return "", err
	}

	defer f.Close()

	if err := f.Chmod(0600); err != nil {
		return "", err
	}

	if _, err := f.Write(data); err != nil {
		return "", err
	}

	return f.Name(), nil
}
