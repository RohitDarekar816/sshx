package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
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

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]
		repoDir := config.GetRepoDir()

		if keyFile == "" {
			return fmt.Errorf("--file is required")
		}

		passphrase, err := requireMFA(repoDir)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := keys.EncryptKeyFile(repoDir, name, keyFile, passphrase); err != nil {
			return fmt.Errorf("encrypting key: %w", err)
		}

		if err := commitAndPush(repoDir, "Add encrypted key "+name); err != nil {
			return fmt.Errorf("git error: %w", err)
		}

		fmt.Println("Key added:", name)
		return nil
	},
}

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List encrypted keys",

	RunE: func(cmd *cobra.Command, args []string) error {

		repoDir := config.GetRepoDir()
		names, err := keys.ListKeys(repoDir)
		if err != nil {
			return fmt.Errorf("listing keys: %w", err)
		}

		if len(names) == 0 {
			fmt.Println("No keys found")
			return nil
		}

		fmt.Println("Keys:")
		for _, n := range names {
			fmt.Printf("- %s\n", n)
		}
		return nil
	},
}

var keyRemoveCmd = &cobra.Command{
	Use:               "remove [name]",
	Short:             "Remove an encrypted key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeKeyNames,

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]
		repoDir := config.GetRepoDir()

		if _, err := requireMFA(repoDir); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := keys.RemoveKey(repoDir, name); err != nil {
			return fmt.Errorf("removing key: %w", err)
		}

		if err := commitAndPush(repoDir, "Remove encrypted key "+name); err != nil {
			return fmt.Errorf("git error: %w", err)
		}

		fmt.Println("Key removed:", name)
		return nil
	},
}

// completeKeyNames provides shell completion for encrypted key names.
func completeKeyNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names, err := keys.ListKeys(config.GetRepoDir())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	keyAddCmd.Flags().StringVar(&keyFile, "file", "", "Path to private key file")

	keyCmd.AddCommand(keyAddCmd)
	keyCmd.AddCommand(keyListCmd)
	keyCmd.AddCommand(keyRemoveCmd)
	rootCmd.AddCommand(keyCmd)
}
