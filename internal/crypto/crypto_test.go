package crypto

import "testing"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plain := "secret-data"
	passphrase := "passphrase"

	enc, err := Encrypt(plain, passphrase)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	dec, err := Decrypt(enc, passphrase)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if dec != plain {
		t.Fatalf("expected %q, got %q", plain, dec)
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	plain := "secret-data"

	enc, err := Encrypt(plain, "correct")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = Decrypt(enc, "wrong")
	if err == nil {
		t.Fatalf("expected decrypt failure with wrong passphrase")
	}
}
