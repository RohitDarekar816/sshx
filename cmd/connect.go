package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/RohitDarekar816/sshx/internal/ssh"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [server]",
	Short: "Connect to server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		repoDir := config.GetRepoDir()

		s, err := server.LoadServer(repoDir, name)
		if err != nil {
			fmt.Println("Server not found")
			return
		}

		ssh.Connect(s)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
