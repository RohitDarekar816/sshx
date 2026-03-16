package keys

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeyPath(t *testing.T) {
	got := KeyPath("/repo", "mykey")
	want := filepath.Join("/repo", "keys", "mykey.yaml")
	if got != want {
		t.Fatalf("KeyPath: got %q, want %q", got, want)
	}
}

func TestSaveAndLoadEncryptedKey(t *testing.T) {
	dir := t.TempDir()

	const encrypted = "encrypted-payload"
	if err := SaveEncryptedKey(dir, "testkey", encrypted); err != nil {
		t.Fatalf("SaveEncryptedKey: %v", err)
	}

	got, err := LoadEncryptedKey(dir, "testkey")
	if err != nil {
		t.Fatalf("LoadEncryptedKey: %v", err)
	}
	if got != encrypted {
		t.Fatalf("LoadEncryptedKey: got %q, want %q", got, encrypted)
	}
}

func TestLoadEncryptedKey_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadEncryptedKey(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing key file")
	}
}

func TestEncryptAndDecryptKey(t *testing.T) {
	dir := t.TempDir()
	passphrase := "test-passphrase"
	plain := "ssh-rsa AAAAB3NzaC1..."

	// Write the plaintext key to a temp file.
	src := filepath.Join(dir, "id_rsa")
	if err := os.WriteFile(src, []byte(plain), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := EncryptKeyFile(dir, "id_rsa", src, passphrase); err != nil {
		t.Fatalf("EncryptKeyFile: %v", err)
	}

	got, err := DecryptKey(dir, "id_rsa", passphrase)
	if err != nil {
		t.Fatalf("DecryptKey: %v", err)
	}
	if string(got) != plain {
		t.Fatalf("DecryptKey: got %q, want %q", got, plain)
	}
}

func TestDecryptKey_WrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "id_rsa")
	if err := os.WriteFile(src, []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := EncryptKeyFile(dir, "id_rsa", src, "correct"); err != nil {
		t.Fatalf("EncryptKeyFile: %v", err)
	}
	_, err := DecryptKey(dir, "id_rsa", "wrong")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestListKeys(t *testing.T) {
	dir := t.TempDir()

	// No keys dir yet – should return empty list.
	names, err := ListKeys(dir)
	if err != nil {
		t.Fatalf("ListKeys (empty): %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(names))
	}

	// Save two keys.
	for _, name := range []string{"alpha", "beta"} {
		if err := SaveEncryptedKey(dir, name, "enc"); err != nil {
			t.Fatalf("SaveEncryptedKey %q: %v", name, err)
		}
	}

	names, err = ListKeys(dir)
	if err != nil {
		t.Fatalf("ListKeys: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(names), names)
	}
}

func TestRemoveKey(t *testing.T) {
	dir := t.TempDir()

	if err := SaveEncryptedKey(dir, "mykey", "enc"); err != nil {
		t.Fatalf("SaveEncryptedKey: %v", err)
	}
	if err := RemoveKey(dir, "mykey"); err != nil {
		t.Fatalf("RemoveKey: %v", err)
	}

	// File should no longer exist.
	if _, err := os.Stat(KeyPath(dir, "mykey")); !os.IsNotExist(err) {
		t.Fatal("expected key file to be removed")
	}
}

func TestWriteTempKey(t *testing.T) {
	data := []byte("key-material")
	path, err := WriteTempKey(data)
	if err != nil {
		t.Fatalf("WriteTempKey: %v", err)
	}
	defer os.Remove(path)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(data) {
		t.Fatalf("WriteTempKey: got %q, want %q", got, data)
	}

	// Permissions should be 0600.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected mode 0600, got %o", info.Mode().Perm())
	}
}
