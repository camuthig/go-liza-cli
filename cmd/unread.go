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
	toUpdate := c.Repositories[pr.Repository.Name].PullRequests[pr.ID]
	toUpdate.MarkUnread(c)
}

func init() {
	rootCmd.AddCommand(unreadCmd)
}
