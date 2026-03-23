package totp

import "testing"

func TestGenerateSecret(t *testing.T) {
	secret, err := GenerateSecret()
	if err != nil {
		t.Fatalf("generate secret failed: %v", err)
	}
	if len(secret) < 16 {
		t.Fatalf("unexpected secret length: %d", len(secret))
	}
}

func TestOTPAuthURL(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	url := OTPAuthURL("sshx", "user@example.com", secret)
	if url == "" {
		t.Fatalf("expected non-empty otpauth url")
	}
}

func TestSecretPath(t *testing.T) {
	path := SecretPath("/home/user/project", "user@example.com")
	expected := "/home/user/project/totp/user_example_com.yaml"
	if path != expected {
		t.Fatalf("expected %s, got %s", expected, path)
	}
}
