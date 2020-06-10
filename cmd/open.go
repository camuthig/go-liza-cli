package cmd

import (
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:       "open",
	ValidArgs: []string{"name", "id"},
	Args:      ValidatePullRequestArgs(),
	Short:     "Open a pull request in your browser and mark it read",
	Run:       RunForPullRequests(openPullRequest),
}

func openPullRequest(c *Config, pr *PullRequest) {
	exec.Command(os.Getenv("BROWSER"), pr.Links.HTML.Href).Run()
	pr.LastRead = time.Now()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
