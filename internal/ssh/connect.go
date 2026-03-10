package ssh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/server"
)

func Connect(s *server.Server, passphrase string) {

	target := fmt.Sprintf("%s@%s", s.User, s.Host)

	sshArgs := []string{}

	if s.Key != "" {
		sshArgs = append(sshArgs, "-i", s.Key)
	}

	if s.Port != 0 {
		sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", s.Port))
	}

	sshArgs = append(sshArgs, target)

	var cmd *exec.Cmd

	if s.Password != "" {
		if passphrase == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter decryption passphrase: ")
			passphrase, _ = reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)
		}

		plainPassword, err := crypto.Decrypt(s.Password, passphrase)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		sshpassArgs := []string{"-p", plainPassword, "ssh"}
		sshpassArgs = append(sshpassArgs, sshArgs...)
		cmd = exec.Command("sshpass", sshpassArgs...)
	} else {
		cmd = exec.Command("ssh", sshArgs...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if s.Password != "" {
			fmt.Println("Connection failed. Make sure 'sshpass' is installed: sudo apt install sshpass")
		} else {
			fmt.Println("Connection failed:", err)
		}
	}
}
