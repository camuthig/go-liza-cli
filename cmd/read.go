package cmd

import (
	"time"

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

func markRead(c *Config, pr *PullRequest) {
	pr.LastRead = time.Now()
}

func init() {
	rootCmd.AddCommand(readCmd)
}
