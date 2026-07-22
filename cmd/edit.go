package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
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
	Use:               "edit [server]",
	Short:             "Edit a server",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeServerNames,

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]
		repoDir := config.GetRepoDir()

		s, err := server.LoadServer(repoDir, name)
		if err != nil {
			return fmt.Errorf("server not found: %s", name)
		}

		if clearPassword && cmd.Flags().Changed("password") {
			return fmt.Errorf("--clear-password cannot be used with --password")
		}

		if clearKeyRef && cmd.Flags().Changed(flagKeyRef) {
			return fmt.Errorf("--clear-key-ref cannot be used with --key-ref")
		}

		if cmd.Flags().Changed("key") && cmd.Flags().Changed(flagKeyRef) {
			return fmt.Errorf("--key and --key-ref cannot be used together")
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

		if cmd.Flags().Changed(flagKeyRef) {
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
				return fmt.Errorf("password cannot be empty; use --clear-password to remove it")
			}

			passphrase, err := readConfirmedSecret("Enter encryption passphrase: ")
			if err != nil {
				return err
			}

			encrypted, err := crypto.Encrypt(editPassword, passphrase)
			if err != nil {
				return fmt.Errorf("encrypting password: %w", err)
			}

			s.Password = encrypted
			s.Key = ""
			s.KeyRef = ""
		}

		if s.Password != "" && (cmd.Flags().Changed("key") || cmd.Flags().Changed(flagKeyRef)) && !cmd.Flags().Changed("password") && !clearPassword {
			fmt.Println("Warning: server has a stored password; key will be ignored unless password is cleared.")
		}

		checked, err := server.ValidateAndRepair(*s)
		if err != nil {
			return fmt.Errorf("invalid server: %w", err)
		}

		if err := server.SaveServer(repoDir, *checked); err != nil {
			return fmt.Errorf("saving server: %w", err)
		}

		if err := commitAndPush(repoDir, "Edit server "+name); err != nil {
			return fmt.Errorf("git error: %w", err)
		}

		fmt.Println("Server updated:", name)
		return nil
	},
}

func init() {

	editCmd.Flags().StringVar(&editHost, "host", "", "Server host")
	editCmd.Flags().StringVar(&editUser, "user", "", "SSH user")
	editCmd.Flags().IntVar(&editPort, "port", 0, "SSH port")
	editCmd.Flags().StringVar(&editKey, "key", "", "SSH key")
	editCmd.Flags().StringVar(&editKeyRef, flagKeyRef, "", "Reference to encrypted key stored in repo")
	editCmd.Flags().StringVar(&editPassword, "password", "", "SSH password")
	editCmd.Flags().BoolVar(&clearPassword, "clear-password", false, "Clear stored password")
	editCmd.Flags().BoolVar(&clearKeyRef, "clear-key-ref", false, "Clear stored key reference")

	rootCmd.AddCommand(editCmd)
}
