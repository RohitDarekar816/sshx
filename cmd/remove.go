package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:               "remove [server]",
	Short:             "Remove a server",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeServerNames,

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]
		repoDir := config.GetRepoDir()

		if err := server.RemoveServer(repoDir, name); err != nil {
			return fmt.Errorf("removing server: %w", err)
		}

		if err := commitAndPush(repoDir, "Remove server "+name); err != nil {
			return fmt.Errorf("git error: %w", err)
		}

		fmt.Println("Server removed:", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
