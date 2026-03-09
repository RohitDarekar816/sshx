package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/git"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var host string
var sshUser string
var port int
var key string
var password string

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		repoDir := config.GetRepoDir()

		s := server.Server{
			Name:     name,
			Host:     host,
			User:     sshUser,
			Port:     port,
			Key:      key,
			Password: password,
		}

		if password != "" && key != "~/.ssh/id_rsa" {
			fmt.Println("Warning: both --key and --password provided. Password will be used.")
		}

		if password != "" {
			s.Key = ""

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter encryption passphrase: ")
			passphrase, _ := reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)

			if passphrase == "" {
				fmt.Println("Error: passphrase cannot be empty")
				return
			}

			fmt.Print("Confirm passphrase: ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)

			if passphrase != confirm {
				fmt.Println("Error: passphrases do not match")
				return
			}

			encrypted, err := crypto.Encrypt(password, passphrase)
			if err != nil {
				fmt.Println("Error encrypting password:", err)
				return
			}

			s.Password = encrypted
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
	addCmd.Flags().StringVar(&sshUser, "user", "root", "SSH user")
	addCmd.Flags().IntVar(&port, "port", 22, "SSH port")
	addCmd.Flags().StringVar(&key, "key", "~/.ssh/id_rsa", "SSH key")
	addCmd.Flags().StringVar(&password, "password", "", "SSH password")

	addCmd.MarkFlagRequired("host")

	rootCmd.AddCommand(addCmd)
}
