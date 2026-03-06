package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers",

	Run: func(cmd *cobra.Command, args []string) {

		repoDir := config.GetRepoDir()

		servers, _ := server.LoadServers(repoDir)

		fmt.Println("Servers:")

		for _, s := range servers {
			fmt.Printf("- %s (%s@%s)\n", s.Name, s.User, s.Host)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
