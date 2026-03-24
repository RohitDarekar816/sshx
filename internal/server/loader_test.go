package server

import "testing"

func TestIsValidServerFile(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{".json", true},
		{".yaml", true},
		{".yml", true},
		{".txt", false},
		{".xml", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := isValidServerFile(tt.ext)
			if result != tt.expected {
				t.Errorf("isValidServerFile(%q) = %v, want %v", tt.ext, result, tt.expected)
			}
		})
	}
}

func TestShouldSkipForDuplicate(t *testing.T) {
	tests := []struct {
		serverName string
		ext        string
		sourceExt  map[string]string
		expected   bool
	}{
		{
			serverName: "server1",
			ext:        ".json",
			sourceExt:  map[string]string{},
			expected:   false,
		},
		{
			serverName: "server1",
			ext:        ".json",
			sourceExt:  map[string]string{"server1": ".yaml"},
			expected:   true,
		},
		{
			serverName: "server1",
			ext:        ".json",
			sourceExt:  map[string]string{"server1": ".yml"},
			expected:   true,
		},
		{
			serverName: "server1",
			ext:        ".yaml",
			sourceExt:  map[string]string{"server1": ".json"},
			expected:   false,
		},
		{
			serverName: "server1",
			ext:        ".json",
			sourceExt:  map[string]string{"server1": ".json"},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.serverName, func(t *testing.T) {
			result := shouldSkipForDuplicate(tt.serverName, tt.ext, tt.sourceExt)
			if result != tt.expected {
				t.Errorf("shouldSkipForDuplicate(%q, %q, ...) = %v, want %v", tt.serverName, tt.ext, result, tt.expected)
			}
		})
	}
}

func TestBuildServerList(t *testing.T) {
	byName := map[string]Server{
		"server1": {Name: "server1", Host: "1.1.1.1"},
		"server2": {Name: "server2", Host: "2.2.2.2"},
	}
	order := []string{"server2", "server1"}

	result := buildServerList(byName, order)

	if len(result) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(result))
	}

	if result[0].Name != "server2" {
		t.Errorf("expected first server to be server2, got %q", result[0].Name)
	}

	if result[1].Name != "server1" {
		t.Errorf("expected second server to be server1, got %q", result[1].Name)
	}
}

func TestBuildServerListEmpty(t *testing.T) {
	byName := map[string]Server{}
	order := []string{}

	result := buildServerList(byName, order)

	if len(result) != 0 {
		t.Fatalf("expected 0 servers, got %d", len(result))
	}
}
