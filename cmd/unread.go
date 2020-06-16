package cmd

import (
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var unreadCmd = &cobra.Command{
	Use:       "unread",
	ValidArgs: []string{"name", "id"},
	Args:      ValidatePullRequestArgs(),
	Short:     "Mark a pull request as read.",
	Run:       RunForPullRequests(unmarkRead),
}

func unmarkRead(c *Config, pr *PullRequestWithRepository) {
	pr.MarkUnread(c)
}

func init() {
	rootCmd.AddCommand(unreadCmd)

	unreadCmd.Flags().BoolVarP(&ForAllPullRequests, "all", "a", false, "Mark all pull requests as unread.")
}
