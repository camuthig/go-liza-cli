package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

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

func openPullRequest(c *Config, pr *PullRequestWithRepository) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", pr.Links.HTML.Href)
	}

	if runtime.GOOS == "linux" {
		cmd = exec.Command(os.Getenv("BROWSER"), pr.Links.HTML.Href)
	}

	if cmd == nil {
		fmt.Println("Unable to open a browser")
		os.Exit(1)
	}

	cmd.Run()

	pr.MarkRead()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
