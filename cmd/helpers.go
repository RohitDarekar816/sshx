package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/keys"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/RohitDarekar816/sshx/internal/ssh"
	"golang.org/x/term"
)

// readLine prints prompt and reads a single trimmed line from stdin.
//
// It reads one byte at a time so it never buffers ahead of the newline, which
// keeps it safe to interleave with readSecret (which reads the raw fd directly).
func readLine(prompt string) (string, error) {
	fmt.Print(prompt)

	var b strings.Builder
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			if buf[0] == '\n' {
				break
			}
			b.WriteByte(buf[0])
		}
		if err != nil {
			if b.Len() == 0 {
				return "", err
			}
			break
		}
	}

	return strings.TrimSpace(strings.TrimRight(b.String(), "\r")), nil
}

// readSecret prints prompt and reads a line without echoing it to the terminal.
//
// When stdin is not a terminal (piped input, tests, automation) it falls back
// to a normal line read so the tool stays scriptable.
func readSecret(prompt string) (string, error) {
	fmt.Print(prompt)

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Not a TTY: fall back to a plain read (prompt already printed).
		line, err := readLine("")
		return line, err
	}

	data, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// readConfirmedSecret reads a secret twice and verifies the two entries match.
func readConfirmedSecret(prompt string) (string, error) {
	first, err := readSecret(prompt)
	if err != nil {
		return "", err
	}
	if first == "" {
		return "", errors.New("passphrase cannot be empty")
	}

	confirm, err := readSecret("Confirm passphrase: ")
	if err != nil {
		return "", err
	}
	if first != confirm {
		return "", errors.New("passphrases do not match")
	}

	return first, nil
}

// commitAndPush commits and pushes repo changes, attributing the commit to the
// locally configured sshx user so the shared history is auditable.
func commitAndPush(repoDir, message string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fall back to an anonymous commit rather than blocking the operation.
		return git.CommitAndPush(repoDir, message, "", "")
	}
	return git.CommitAndPush(repoDir, message, cfg.UserName, cfg.UserEmail)
}

// connectToServer resolves a server profile, authenticates the user, materializes
// any encrypted key, and hands off to the system ssh client. It is shared by the
// `connect` command and the bare `sshx <server>` shortcut.
func connectToServer(name string) error {
	repoDir := config.GetRepoDir()

	s, err := server.LoadServer(repoDir, name)
	if err != nil {
		return fmt.Errorf("server not found: %s", name)
	}

	passphrase, err := requireMFA(repoDir)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if s.Password != "" {
		s.Key = ""
		s.KeyRef = ""
	}

	if s.KeyRef != "" {
		keyData, err := keys.DecryptKey(repoDir, s.KeyRef, passphrase)
		if err != nil {
			return fmt.Errorf("decrypting key: %w", err)
		}

		tempKey, err := keys.WriteTempKey(keyData)
		if err != nil {
			return fmt.Errorf("writing temp key: %w", err)
		}
		defer os.Remove(tempKey)

		s.Key = tempKey
	}

	return ssh.Connect(s, passphrase)
}
