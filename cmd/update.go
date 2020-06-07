package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull the latest pull request information from BitBucket.",
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()
		updateWatchedPullrequests(&c)
		updatePullRequestActivity(&c)
		c.Write()
		fmt.Println("Update completed")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateWatchedPullrequests(c *Config) {
	for n, r := range c.Repositories {
		latest := make(map[int]*PullRequest)
		prs := GetWatchedPullRequests(*c, n)

		for _, pr := range prs {
			e, found := r.PullRequests[pr.ID]

			if found {
				latest[e.ID] = e
			} else {
				latest[pr.ID] = pr
			}
		}

		r.PullRequests = latest
		c.Repositories[n] = r
	}

}

func updatePullRequestActivity(c *Config) {
	var next string

	for n, r := range c.Repositories {
		for _, pr := range r.PullRequests {
			var updates []PullRequestUpdate
			var until time.Time
			if pr.LastRead.Before(pr.LastUpdated) {
				updates = pr.Updates
				until = pr.LastUpdated
			} else {
				updates = make([]PullRequestUpdate, 0)
				until = pr.LastRead
			}
			for hasNext := true; hasNext; hasNext = (next != "") {
				var us []PullRequestUpdate
				us, next = GetPullRequestActivity(*c, n, pr.ID, next)

				for _, u := range us {
					// Break the loop once it is older than the last user read
					if u.Date.Before(until) {
						next = ""
						break
					}

					// Ignore updates by this user
					if u.Author.UUID == c.UserUUID {
						continue
					}

					// Ignore approvals if this user is only a reviewer
					if pr.Author.UUID != c.UserUUID && u.ActivityType == Approval {
						continue
					}

					updates = append(updates, u)
				}
			}

			pr.Updates = updates
			pr.LastUpdated = time.Now()
		}
	}
}
