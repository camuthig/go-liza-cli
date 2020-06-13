package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var countOnly bool

// updatesCmd represents the updates command
var updatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()

		if countOnly {
			count := 0
			for _, r := range c.Repositories {
				for _, pr := range r.PullRequests {
					count += pr.UnreadUpdatesCount
				}
			}

			fmt.Println(count)
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Repository", "Pull Request", "# of Updates", "Last Updated", "Last Read", "Link"})

		for _, r := range c.Repositories {
			for _, pr := range r.PullRequests {
				title := pr.Title
				if len(title) > 35 {
					title = title[0:35] + "..."
				}
				t.AppendRow(table.Row{r.Name, title, pr.UnreadUpdatesCount, pr.LastUpdated, pr.ReadAt, pr.Links.HTML.Href})
			}
		}

		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(updatesCmd)

	updatesCmd.Flags().BoolVarP(&countOnly, "count", "c", false, "Only return the count of updates.")
}
