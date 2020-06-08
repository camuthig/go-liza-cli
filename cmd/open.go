package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open {name} {id}",
	Args:  cobra.ExactArgs(2),
	Short: "Open a pull request in your browser and mark it read",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		id, err := strconv.Atoi(args[1])

		if err != nil {
			fmt.Printf("Invalid ID %s", args[1])
			os.Exit(1)
		}

		c := ParseConfig()

		r, found := c.Repositories[name]

		if !found {
			fmt.Printf("Not watching repository %s", name)
			os.Exit(1)
		}

		pr, found := r.PullRequests[id]

		if !found {
			fmt.Printf("Pull request %d not found", id)
			os.Exit(1)
		}

		exec.Command(os.Getenv("BROWSER"), pr.Links.HTML.Href).Run()
		pr.LastRead = time.Now()

		c.Write()
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
