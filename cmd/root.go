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

var rootCmd = &cobra.Command{
	Use:   "sshx",
	Short: "SSH profile manager",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			return
		}

		serverName := args[0]
		repoDir := config.GetRepoDir()

		s, err := server.LoadServer(repoDir, serverName)

		if err != nil {
			fmt.Println("Server not found:", serverName)
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

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
