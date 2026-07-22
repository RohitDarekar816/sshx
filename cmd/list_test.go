package cmd

import (
	"testing"

	"github.com/RohitDarekar816/sshx/internal/server"
)

func TestAuthKind(t *testing.T) {
	cases := []struct {
		name string
		in   server.Server
		want string
	}{
		{"password wins", server.Server{Password: "x", Key: "k", KeyRef: "r"}, "password"},
		{"key-ref", server.Server{KeyRef: "prod"}, "key-ref:prod"},
		{"inline key", server.Server{Key: "~/.ssh/id_rsa"}, "key"},
		{"none", server.Server{}, "-"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := authKind(c.in); got != c.want {
				t.Errorf("authKind() = %q, want %q", got, c.want)
			}
		})
	}
}
