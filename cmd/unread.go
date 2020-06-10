package cmd

import (
	"time"

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

func unmarkRead(c *Config, pr *PullRequest) {
	pr.LastRead = pr.LastUpdated.Add(-1 * time.Second)
}

func init() {
	rootCmd.AddCommand(unreadCmd)
}
