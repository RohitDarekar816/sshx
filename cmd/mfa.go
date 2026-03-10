package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/totp"
)

func requireMFA(repoDir string) (string, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	email := strings.TrimSpace(cfg.UserEmail)
	reader := bufio.NewReader(os.Stdin)

	if email == "" {
		fmt.Print("Enter your email: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
		if email == "" {
			return "", errors.New("email is required")
		}
		cfg.UserEmail = email
		_ = config.SaveConfig(cfg)
	}

	if !totp.HasSecret(repoDir, email) {
		return "", errors.New("TOTP not set up. Run `sshx auth` first")
	}

	fmt.Print("Enter vault passphrase: ")
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)
	if passphrase == "" {
		return "", errors.New("passphrase cannot be empty")
	}

	encrypted, err := totp.LoadEncryptedSecret(repoDir, email)
	if err != nil {
		return "", err
	}

	secret, err := crypto.Decrypt(encrypted, passphrase)
	if err != nil {
		return "", errors.New("invalid passphrase or corrupt TOTP secret")
	}

	fmt.Print("Enter TOTP code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	if !totp.Validate(secret, code, time.Now()) {
		return "", errors.New("invalid TOTP code")
	}

	return passphrase, nil
}
