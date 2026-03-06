package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var host string
var user string
var port int
var key string

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		repoDir := config.GetRepoDir()

		s := server.Server{
			Name: name,
			Host: host,
			User: user,
			Port: port,
			Key:  key,
		}

		err := server.SaveServer(repoDir, s)
		if err != nil {
			fmt.Println("Error saving server:", err)
			return
		}

		err = git.CommitAndPush(repoDir, "Add server "+name)
		if err != nil {
			fmt.Println("Git error:", err)
			return
		}

		fmt.Println("Server added:", name)
	},
}

func init() {

	addCmd.Flags().StringVar(&host, "host", "", "Server host")
	addCmd.Flags().StringVar(&user, "user", "root", "SSH user")
	addCmd.Flags().IntVar(&port, "port", 22, "SSH port")
	addCmd.Flags().StringVar(&key, "key", "~/.ssh/id_rsa", "SSH key")

	addCmd.MarkFlagRequired("host")

	rootCmd.AddCommand(addCmd)
}
