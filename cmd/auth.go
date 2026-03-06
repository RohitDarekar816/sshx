package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/user"
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

		reader := bufio.NewReader(os.Stdin)

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

		fmt.Println("sshx authentication complete")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
