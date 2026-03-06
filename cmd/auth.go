package cmd

import (
	"fmt"
	"os"

	package cmd

import (
	"fmt"
	"os"

	"github.com/rohit/sshx/internal/config"
	"github.com/rohit/sshx/internal/git"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [repo-url]",
	Short: "Authenticate sshx with a git repository",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		repoURL := args[0]

		fmt.Println("Initializing sshx...")

		// create ~/.sshx
		err := config.InitBaseDir()
		if err != nil {
			fmt.Println("Error creating sshx directory:", err)
			os.Exit(1)
		}

		repoDir := config.GetRepoDir()

		// clone repo
		err = git.CloneRepo(repoURL, repoDir)
		if err != nil {
			fmt.Println("Error cloning repo:", err)
			os.Exit(1)
		}

		fmt.Println("sshx authenticated successfully!")
		fmt.Println("Repo:", repoURL)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
)

var authCmd = &cobra.Command{
	Use:   "auth [repo-url]",
	Short: "Authenticate sshx with a git repository",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		repoURL := args[0]

		fmt.Println("Initializing sshx...")

		// create ~/.sshx
		err := config.InitBaseDir()
		if err != nil {
			fmt.Println("Error creating sshx directory:", err)
			os.Exit(1)
		}

		repoDir := config.GetRepoDir()

		// clone repo
		err = git.CloneRepo(repoURL, repoDir)
		if err != nil {
			fmt.Println("Error cloning repo:", err)
			os.Exit(1)
		}

		fmt.Println("sshx authenticated successfully!")
		fmt.Println("Repo:", repoURL)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
