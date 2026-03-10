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

var editHost string
var editUser string
var editPort int
var editKey string
var editKeyRef string
var editPassword string
var clearPassword bool
var clearKeyRef bool

var editCmd = &cobra.Command{
	Use:   "edit [server]",
	Short: "Edit a server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		repoDir := config.GetRepoDir()

		s, err := server.LoadServer(repoDir, name)
		if err != nil {
			fmt.Println("Server not found:", name)
			return
		}

		if clearPassword && cmd.Flags().Changed("password") {
			fmt.Println("Error: --clear-password cannot be used with --password")
			return
		}

		if clearKeyRef && cmd.Flags().Changed("key-ref") {
			fmt.Println("Error: --clear-key-ref cannot be used with --key-ref")
			return
		}

		if cmd.Flags().Changed("key") && cmd.Flags().Changed("key-ref") {
			fmt.Println("Error: --key and --key-ref cannot be used together")
			return
		}

		if cmd.Flags().Changed("host") {
			s.Host = editHost
		}

		if cmd.Flags().Changed("user") {
			s.User = editUser
		}

		if cmd.Flags().Changed("port") {
			s.Port = editPort
		}

		if cmd.Flags().Changed("key") {
			s.Key = editKey
		}

		if cmd.Flags().Changed("key-ref") {
			s.KeyRef = editKeyRef
			s.Key = ""
		}

		if clearKeyRef {
			s.KeyRef = ""
		}

		if clearPassword {
			s.Password = ""
		}

		if cmd.Flags().Changed("password") {
			if editPassword == "" {
				fmt.Println("Error: password cannot be empty. Use --clear-password to remove it.")
				return
			}

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

			encrypted, err := crypto.Encrypt(editPassword, passphrase)
			if err != nil {
				fmt.Println("Error encrypting password:", err)
				return
			}

			s.Password = encrypted
			s.Key = ""
			s.KeyRef = ""
		}

		if s.Password != "" && (cmd.Flags().Changed("key") || cmd.Flags().Changed("key-ref")) && !cmd.Flags().Changed("password") && !clearPassword {
			fmt.Println("Warning: server has a stored password; key will be ignored unless password is cleared.")
		}

		if err := server.SaveServer(repoDir, *s); err != nil {
			fmt.Println("Error saving server:", err)
			return
		}

		if err := git.CommitAndPush(repoDir, "Edit server "+name); err != nil {
			fmt.Println("Git error:", err)
			return
		}

		fmt.Println("Server updated:", name)
	},
}

func init() {

	editCmd.Flags().StringVar(&editHost, "host", "", "Server host")
	editCmd.Flags().StringVar(&editUser, "user", "", "SSH user")
	editCmd.Flags().IntVar(&editPort, "port", 0, "SSH port")
	editCmd.Flags().StringVar(&editKey, "key", "", "SSH key")
	editCmd.Flags().StringVar(&editKeyRef, "key-ref", "", "Reference to encrypted key stored in repo")
	editCmd.Flags().StringVar(&editPassword, "password", "", "SSH password")
	editCmd.Flags().BoolVar(&clearPassword, "clear-password", false, "Clear stored password")
	editCmd.Flags().BoolVar(&clearKeyRef, "clear-key-ref", false, "Clear stored key reference")

	rootCmd.AddCommand(editCmd)
}
