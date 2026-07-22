package cmd

import (
	"github.com/spf13/cobra"
)

// version is overridable at build time via -ldflags "-X ...cmd.version=X.Y.Z".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "sshx [server]",
	Short:         "Git-backed SSH profile manager",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MaximumNArgs(1),
	// Bare `sshx <server>` is a shortcut for `sshx connect <server>`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return connectToServer(args[0])
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
