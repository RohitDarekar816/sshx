package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/RohitDarekar816/sshx/internal/ssh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshx",
	Short: "SSH profile manager",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			return
		}

		serverName := args[0]

		serverPath := filepath.Join(
			config.GetRepoDir(),
			"servers",
			serverName+".yaml",
		)

		s, err := server.LoadServer(serverPath)

		if err != nil {
			fmt.Println("Server not found:", serverName)
			return
		}

		ssh.Connect(s)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
