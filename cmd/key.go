package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/keys"
	"github.com/spf13/cobra"
)

var keyFile string

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage encrypted SSH keys",
}

var keyAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add an encrypted key",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		repoDir := config.GetRepoDir()

		if keyFile == "" {
			fmt.Println("Error: --file is required")
			return
		}

		passphrase, err := requireMFA(repoDir)
		if err != nil {
			fmt.Println("Authentication failed:", err)
			return
		}

		if err := keys.EncryptKeyFile(repoDir, name, keyFile, passphrase); err != nil {
			fmt.Println("Error encrypting key:", err)
			return
		}

		if err := git.CommitAndPush(repoDir, "Add encrypted key "+name); err != nil {
			fmt.Println("Git error:", err)
			return
		}

		fmt.Println("Key added:", name)
	},
}

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List encrypted keys",

	Run: func(cmd *cobra.Command, args []string) {

		repoDir := config.GetRepoDir()
		names, err := keys.ListKeys(repoDir)
		if err != nil {
			fmt.Println("Error listing keys:", err)
			return
		}

		if len(names) == 0 {
			fmt.Println("No keys found")
			return
		}

		fmt.Println("Keys:")
		for _, n := range names {
			fmt.Printf("- %s\n", n)
		}
	},
}

var keyRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an encrypted key",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		repoDir := config.GetRepoDir()

		_, err := requireMFA(repoDir)
		if err != nil {
			fmt.Println("Authentication failed:", err)
			return
		}

		if err := keys.RemoveKey(repoDir, name); err != nil {
			fmt.Println("Error removing key:", err)
			return
		}

		if err := git.CommitAndPush(repoDir, "Remove encrypted key "+name); err != nil {
			fmt.Println("Git error:", err)
			return
		}

		fmt.Println("Key removed:", name)
	},
}

func init() {
	keyAddCmd.Flags().StringVar(&keyFile, "file", "", "Path to private key file")

	keyCmd.AddCommand(keyAddCmd)
	keyCmd.AddCommand(keyListCmd)
	keyCmd.AddCommand(keyRemoveCmd)
	rootCmd.AddCommand(keyCmd)
}
