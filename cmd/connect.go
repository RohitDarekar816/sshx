package cmd

import (
	"fmt"
	"os"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/keys"
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

		passphrase, err := requireMFA(repoDir)
		if err != nil {
			fmt.Println("Authentication failed:", err)
			return
		}

		if s.Password != "" {
			s.Key = ""
			s.KeyRef = ""
		}

		var tempKey string
		if s.KeyRef != "" {
			keyData, err := keys.DecryptKey(repoDir, s.KeyRef, passphrase)
			if err != nil {
				fmt.Println("Error decrypting key:", err)
				return
			}

			tempKey, err = keys.WriteTempKey(keyData)
			if err != nil {
				fmt.Println("Error writing temp key:", err)
				return
			}
			defer os.Remove(tempKey)

			s.Key = tempKey
		}

		ssh.Connect(s, passphrase)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
