package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [server]",
	Short: "Remove a server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		repoDir := config.GetRepoDir()

		if err := server.RemoveServer(repoDir, name); err != nil {
			fmt.Println("Error removing server:", err)
			return
		}

		if err := git.CommitAndPush(repoDir, "Remove server "+name); err != nil {
			fmt.Println("Git error:", err)
			return
		}

		fmt.Println("Server removed:", name)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
