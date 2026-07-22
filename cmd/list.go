package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/RohitDarekar816/sshx/internal/config"
	"github.com/RohitDarekar816/sshx/internal/server"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers",

	RunE: func(cmd *cobra.Command, args []string) error {

		repoDir := config.GetRepoDir()

		servers, err := server.LoadServers(repoDir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no profiles found. Run `sshx auth` first")
			}
			return fmt.Errorf("loading servers: %w", err)
		}

		if len(servers) == 0 {
			fmt.Println("No servers configured. Add one with `sshx add <name> --host <host>`.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTARGET\tPORT\tAUTH")
		for _, s := range servers {
			fmt.Fprintf(w, "%s\t%s@%s\t%d\t%s\n", s.Name, s.User, s.Host, s.Port, authKind(s))
		}
		return w.Flush()
	},
}

// authKind reports how a profile authenticates, for display in `list`.
func authKind(s server.Server) string {
	switch {
	case s.Password != "":
		return "password"
	case s.KeyRef != "":
		return "key-ref:" + s.KeyRef
	case s.Key != "":
		return "key"
	default:
		return "-"
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
