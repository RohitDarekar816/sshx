package cmd

import (
	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

// completeServerNames provides shell completion for commands that take a server
// name argument. It returns the profile names known in the local repo.
func completeServerNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	servers, err := server.LoadServers(config.GetRepoDir())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	names := make([]string, 0, len(servers))
	for _, s := range servers {
		names = append(names, s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
