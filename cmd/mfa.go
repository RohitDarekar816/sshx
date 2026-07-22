package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/totp"
	"github.com/mdp/qrterminal/v3"
)

// envPassphrase lets automation supply the vault passphrase non-interactively.
const envPassphrase = "SSHX_PASSPHRASE"

func requireMFA(repoDir string) (string, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	email := strings.TrimSpace(cfg.UserEmail)

	if email == "" {
		email, err = readLine("Enter your email: ")
		if err != nil {
			return "", err
		}
		if email == "" {
			return "", errors.New("email is required")
		}
		cfg.UserEmail = email
		_ = config.SaveConfig(cfg)
	}

	if !totp.HasSecret(repoDir, email) {
		return "", errors.New("TOTP not set up. Run `sshx auth` first")
	}

	passphrase := os.Getenv(envPassphrase)
	if passphrase == "" {
		passphrase, err = readSecret("Enter vault passphrase: ")
		if err != nil {
			return "", err
		}
	}
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

	code, err := readLine("Enter TOTP code: ")
	if err != nil {
		return "", err
	}

	if !totp.Validate(secret, code, time.Now()) {
		return "", errors.New("invalid TOTP code")
	}

	return passphrase, nil
}

func setupTOTP(repoDir, email string) error {

	secret, err := totp.GenerateSecret()
	if err != nil {
		return err
	}

	fmt.Println("TOTP setup required.")
	fmt.Println("Scan this QR code with your authenticator app:")
	otpURL := totp.OTPAuthURL("sshx", email, secret)
	qrterminal.GenerateHalfBlock(otpURL, qrterminal.L, os.Stdout)
	fmt.Println("If you cannot scan the QR code, use this URL:")
	fmt.Println(otpURL)
	fmt.Println("Manual secret:", secret)

	code, err := readLine("Enter the current TOTP code to verify: ")
	if err != nil {
		return err
	}
	if !totp.Validate(secret, code, time.Now()) {
		return fmt.Errorf("invalid TOTP code")
	}

	passphrase, err := readConfirmedSecret("Create vault passphrase: ")
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(secret, passphrase)
	if err != nil {
		return err
	}

	return totp.SaveEncryptedSecret(repoDir, email, encrypted)
}
