package cmd

import (
	"fmt"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/crypto"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var host string
var sshUser string
var port int
var key string
var keyRef string
var password string

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new server",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]

		repoDir := config.GetRepoDir()

		s := server.Server{
			Name:     name,
			Host:     host,
			User:     sshUser,
			Port:     port,
			Key:      key,
			KeyRef:   keyRef,
			Password: password,
		}

		if keyRef != "" && cmd.Flags().Changed("key") {
			return fmt.Errorf("--key and --key-ref cannot be used together")
		}

		if keyRef != "" {
			s.Key = ""
		}

		if password != "" && (cmd.Flags().Changed("key") || keyRef != "") {
			fmt.Println("Warning: password provided. SSH key will be ignored.")
		}

		if password != "" {
			s.Key = ""

			passphrase, err := readConfirmedSecret("Enter encryption passphrase: ")
			if err != nil {
				return err
			}

			encrypted, err := crypto.Encrypt(password, passphrase)
			if err != nil {
				return fmt.Errorf("encrypting password: %w", err)
			}

			s.Password = encrypted
		}

		checked, err := server.ValidateAndRepair(s)
		if err != nil {
			return fmt.Errorf("invalid server: %w", err)
		}

		if err := server.SaveServer(repoDir, *checked); err != nil {
			return fmt.Errorf("saving server: %w", err)
		}

		if err := commitAndPush(repoDir, "Add server "+name); err != nil {
			return fmt.Errorf("git error: %w", err)
		}

		fmt.Println("Server added:", name)
		return nil
	},
}

func init() {

	addCmd.Flags().StringVar(&host, "host", "", "Server host")
	addCmd.Flags().StringVar(&sshUser, "user", "root", "SSH user")
	addCmd.Flags().IntVar(&port, "port", 22, "SSH port")
	addCmd.Flags().StringVar(&key, "key", "~/.ssh/id_rsa", "SSH key")
	addCmd.Flags().StringVar(&keyRef, flagKeyRef, "", "Reference to encrypted key stored in repo")
	addCmd.Flags().StringVar(&password, "password", "", "SSH password")

	addCmd.MarkFlagRequired("host")

	rootCmd.AddCommand(addCmd)
}
