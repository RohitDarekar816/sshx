package cmd

import (
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:               "connect [server]",
	Short:             "Connect to a server",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeServerNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		return connectToServer(args[0])
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
