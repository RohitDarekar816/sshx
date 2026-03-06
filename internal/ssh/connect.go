package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/RohitDarekar816/sshx/internal/server"
)

func Connect(s *server.Server) {

	target := fmt.Sprintf("%s@%s", s.User, s.Host)

	args := []string{}

	if s.Key != "" {
		args = append(args, "-i", s.Key)
	}

	if s.Port != 0 {
		args = append(args, "-p", fmt.Sprintf("%d", s.Port))
	}

	args = append(args, target)

	cmd := exec.Command("ssh", args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Run()
}
