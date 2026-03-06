package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
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

		address := fmt.Sprintf("%s@%s", s.User, s.Host)

		c := exec.Command("ssh", address)

		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
