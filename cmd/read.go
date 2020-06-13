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
	toUpdate := c.Repositories[pr.Repository.Name].PullRequests[pr.ID]
	toUpdate.MarkRead()
}

func init() {
	rootCmd.AddCommand(readCmd)
}
