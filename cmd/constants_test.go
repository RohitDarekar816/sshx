package cmd

import (
	"testing"
)

func TestFlagKeyRefConstant(t *testing.T) {
	if flagKeyRef != "key-ref" {
		t.Errorf("flagKeyRef = %q, want %q", flagKeyRef, "key-ref")
	}
}

func TestEditCmdHasKeyRefFlag(t *testing.T) {
	f := editCmd.Flags().Lookup(flagKeyRef)
	if f == nil {
		t.Errorf("editCmd does not have flag %q", flagKeyRef)
	}
}

func TestAddCmdHasKeyRefFlag(t *testing.T) {
	f := addCmd.Flags().Lookup(flagKeyRef)
	if f == nil {
		t.Errorf("addCmd does not have flag %q", flagKeyRef)
	}
}
