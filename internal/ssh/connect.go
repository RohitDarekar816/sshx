package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/server"
	"golang.org/x/term"
)

// Connect builds the ssh invocation for the given profile and hands control to
// the system ssh client, forwarding stdio. It returns an error if the profile's
// password cannot be decrypted or the ssh process exits non-zero.
func Connect(s *server.Server, passphrase string) error {

	target := fmt.Sprintf("%s@%s", s.User, s.Host)

	sshArgs := []string{}

	if s.Key != "" {
		sshArgs = append(sshArgs, "-i", s.Key)
	}

	if s.Port != 0 {
		sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", s.Port))
	}

	if s.Password != "" {
		sshArgs = append(sshArgs,
			"-o", "PreferredAuthentications=password,keyboard-interactive",
			"-o", "PubkeyAuthentication=no",
			"-o", "PasswordAuthentication=yes",
		)
	}

	sshArgs = append(sshArgs, target)

	var cmd *exec.Cmd
	usedSshpass := false

	if s.Password != "" {
		if passphrase == "" {
			fmt.Print("Enter decryption passphrase: ")
			if term.IsTerminal(int(os.Stdin.Fd())) {
				data, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				if err != nil {
					return fmt.Errorf("reading passphrase: %w", err)
				}
				passphrase = string(data)
			} else {
				fmt.Fscanln(os.Stdin, &passphrase)
			}
		}

		plainPassword, err := crypto.Decrypt(s.Password, passphrase)
		if err != nil {
			return fmt.Errorf("decrypting password: %w", err)
		}

		if _, err := exec.LookPath("sshpass"); err == nil {
			sshpassArgs := []string{"-p", plainPassword, "ssh"}
			sshpassArgs = append(sshpassArgs, sshArgs...)
			cmd = exec.Command("sshpass", sshpassArgs...)
			usedSshpass = true
		} else {
			fmt.Println("sshpass not found. Falling back to interactive SSH password prompt.")
			cmd = exec.Command("ssh", sshArgs...)
		}
	} else {
		cmd = exec.Command("ssh", sshArgs...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if usedSshpass {
			return fmt.Errorf("connection failed: %w\nsshpass was detected and used; verify the password, server auth settings, and SSH reachability", err)
		}
		return fmt.Errorf("connection failed: %w", err)
	}

	return nil
}
