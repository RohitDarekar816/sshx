package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/totp"
	"github.com/RohitDarekar816/sshx/internal/user"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [repo-url]",
	Short: "Authenticate sshx with a git repository",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		repoURL := args[0]

		fmt.Println("Initializing sshx...")

		if err := config.InitBaseDir(); err != nil {
			return fmt.Errorf("initializing sshx directory: %w", err)
		}

		confirm, err := readLine("Confirm repository is private (yes/no): ")
		if err != nil {
			return err
		}
		confirm = strings.ToLower(confirm)
		if confirm != "yes" && confirm != "y" {
			return fmt.Errorf("sshx requires a private repository")
		}

		repoDir := config.GetRepoDir()

		fmt.Println("Cloning repository:", repoURL)

		if err := git.CloneOrInitRepo(repoURL, repoDir); err != nil {
			return fmt.Errorf("preparing repo: %w", err)
		}

		usersFilePath := filepath.Join(repoDir, "users.json")

		users, err := user.LoadUsers(usersFilePath)
		if err != nil {
			return fmt.Errorf("reading users: %w", err)
		}

		name, err := readLine("Enter your name: ")
		if err != nil {
			return err
		}
		email, err := readLine("Enter your email: ")
		if err != nil {
			return err
		}

		// Persist local config early so commitAndPush attributes commits correctly.
		cfg := &config.Config{
			RepoURL:   repoURL,
			RepoDir:   repoDir,
			UserName:  name,
			UserEmail: email,
		}
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("saving local config: %w", err)
		}

		if user.UserExists(users, email) {
			fmt.Println("User already registered")
		} else {
			user.AddUser(users, name, email)

			if err := user.SaveUsers(usersFilePath, users); err != nil {
				return fmt.Errorf("saving users: %w", err)
			}

			if err := commitAndPush(repoDir, "Add new sshx user"); err != nil {
				return fmt.Errorf("pushing changes: %w", err)
			}

			fmt.Println("User registered successfully")
		}

		if !totp.HasSecret(repoDir, email) {
			if err := setupTOTP(repoDir, email); err != nil {
				return fmt.Errorf("TOTP setup failed: %w", err)
			}

			if err := commitAndPush(repoDir, "Add TOTP secret for "+email); err != nil {
				return fmt.Errorf("pushing changes: %w", err)
			}
		}

		fmt.Println("sshx authentication complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
