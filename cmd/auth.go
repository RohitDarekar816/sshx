package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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

		err := git.CloneRepo(repoURL, repoDir)
		if err != nil {
			fmt.Println("Error cloning repo:", err)
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

		name = name[:len(name)-1]
		email = email[:len(email)-1]

		if user.UserExists(users, email) {

			fmt.Println("User already registered")

		} else {

			user.AddUser(users, name, email)

			user.SaveUsers(usersFilePath, users)

			git.CommitAndPush(repoDir, "Add new sshx user")

			fmt.Println("User registered successfully")
		}

		fmt.Println("sshx authentication complete")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
