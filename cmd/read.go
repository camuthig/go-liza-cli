package cmd

import (
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var readCmd = &cobra.Command{
	Use:       "read",
	ValidArgs: []string{"name", "id"},
	Args:      ValidatePullRequestArgs(),
	Short:     "Mark a pull request as read.",
	Run:       RunForPullRequests(markRead),
}

func markRead(c *Config, pr *PullRequestWithRepository) {
	pr.MarkRead()
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.Flags().BoolVarP(&ForAllPullRequests, "all", "a", false, "Mark all pull requests as read.")
}
