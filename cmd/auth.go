package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/totp"
	"github.com/RohitDarekar816/sshx/internal/user"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [repo-url]",
	Short: "Authenticate sshx with a git repository",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		repoURL := args[0]

		fmt.Println("Initializing sshx...")

		config.InitBaseDir()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Confirm repository is private (yes/no): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))
		if confirm != "yes" && confirm != "y" {
			fmt.Println("sshx requires a private repository")
			return
		}

		repoDir := config.GetRepoDir()

		fmt.Println("Cloning repository:", repoURL)

		err := git.CloneOrInitRepo(repoURL, repoDir)
		if err != nil {
			fmt.Println("Error preparing repo:", err)
			return
		}

		usersFilePath := filepath.Join(repoDir, "users.json")

		users, err := user.LoadUsers(usersFilePath)
		if err != nil {
			fmt.Println("Error reading users:", err)
			return
		}

		fmt.Print("Enter your name: ")
		name, _ := reader.ReadString('\n')

		fmt.Print("Enter your email: ")
		email, _ := reader.ReadString('\n')

		name = strings.TrimSpace(name)
		email = strings.TrimSpace(email)

		if user.UserExists(users, email) {

			fmt.Println("User already registered")

		} else {

			user.AddUser(users, name, email)

			err = user.SaveUsers(usersFilePath, users)
			if err != nil {
				fmt.Println("Error saving users:", err)
				return
			}

			err = git.CommitAndPush(repoDir, "Add new sshx user")
			if err != nil {
				fmt.Println("Error pushing changes:", err)
				return
			}

			fmt.Println("User registered successfully")
		}

		cfg := &config.Config{
			RepoURL:   repoURL,
			RepoDir:   repoDir,
			UserName:  name,
			UserEmail: email,
		}
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving local config:", err)
			return
		}

		if !totp.HasSecret(repoDir, email) {
			if err := setupTOTP(repoDir, email, reader); err != nil {
				fmt.Println("TOTP setup failed:", err)
				return
			}

			if err := git.CommitAndPush(repoDir, "Add TOTP secret for "+email); err != nil {
				fmt.Println("Error pushing changes:", err)
				return
			}
		}

		fmt.Println("sshx authentication complete")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func setupTOTP(repoDir, email string, reader *bufio.Reader) error {

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

	fmt.Print("Enter the current TOTP code to verify: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)
	if !totp.Validate(secret, code, time.Now()) {
		return fmt.Errorf("invalid TOTP code")
	}

	fmt.Print("Create vault passphrase: ")
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	fmt.Print("Confirm passphrase: ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != passphrase {
		return fmt.Errorf("passphrases do not match")
	}

	encrypted, err := crypto.Encrypt(secret, passphrase)
	if err != nil {
		return err
	}

	return totp.SaveEncryptedSecret(repoDir, email, encrypted)
}
