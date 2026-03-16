package server

import "testing"

const testHost = "1.2.3.4"

func TestValidateAndRepairDefaults(t *testing.T) {
	input := Server{
		Name: "prod",
		Host: testHost,
	}

	out, err := ValidateAndRepair(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.User != "root" {
		t.Fatalf("expected default user root, got %q", out.User)
	}

	if out.Port != 22 {
		t.Fatalf("expected default port 22, got %d", out.Port)
	}

	if out.Key != "~/.ssh/id_rsa" {
		t.Fatalf("expected default key, got %q", out.Key)
	}
}

func TestValidateAndRepairInvalidPort(t *testing.T) {
	input := Server{
		Name: "prod",
		Host: testHost,
		Port: 70000,
	}

	_, err := ValidateAndRepair(input)
	if err == nil {
		t.Fatalf("expected error for invalid port")
	}
}

func TestValidateAndRepairAuthPrecedence(t *testing.T) {
	input := Server{
		Name:     "prod",
		Host:     testHost,
		User:     "ubuntu",
		Port:     22,
		Key:      "~/.ssh/id_rsa",
		KeyRef:   "prod",
		Password: "enc",
	}

	out, err := ValidateAndRepair(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Password == "" {
		t.Fatalf("expected password to remain")
	}

	if out.Key != "" || out.KeyRef != "" {
		t.Fatalf("expected key and key_ref to be cleared when password is set")
	}
}
