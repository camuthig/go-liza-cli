package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
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

		text.EnableColors()

		for _, r := range c.Repositories {
			fmt.Println(text.Bold.Sprint(r.Name))

			for _, pr := range r.PullRequests {
				if pr.UnreadUpdatesCount == 0 {
					continue
				}

				approvals := 0
				comments := 0
				updates := 0
				for _, u := range pr.Updates {
					if !u.IsNewToUserUUID(c.UserUUID, pr.ReadAt) {
						continue
					}

					if u.ActivityType == Update {
						updates++
					} else if u.ActivityType == Approval {
						approvals++
					} else if u.ActivityType == Comment {
						comments++
					}
				}

				fmt.Printf("    [%d] %s\n", pr.ID, pr.Title)
				fmt.Printf("        %s\n", pr.Links.HTML.Href)
				fmt.Println(text.Colors{text.Bold}.Sprintf("        Updates: %d", updates))
				fmt.Println(text.Colors{text.Bold, text.FgGreen}.Sprintf("        Approvals: %d", approvals))
				fmt.Println(text.Colors{text.Bold, text.FgBlue}.Sprintf("        Comments: %d", comments))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updatesCmd)

	updatesCmd.Flags().BoolVarP(&countOnly, "count", "c", false, "Only return the count of updates.")
}
