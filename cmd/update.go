package cmd

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull the latest pull request information from BitBucket.",
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()
		updateWatchedPullrequests(&c)
		updatePullRequestsActivity(&c)
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

	err := beeep.Notify("Updated pull requests", "Do you hear me?!?", "assets/bitbucket.png")
	if err != nil {
		panic(err)
	}

}

func updatePullRequestsActivity(c *Config) {
	for _, r := range c.Repositories {
		for _, pr := range r.PullRequests {
			updatePullRequestActivity(c, r, pr)
		}
	}
}

func updatePullRequestActivity(c *Config, r *Repository, pr *PullRequest) {
	var next string
	var latest []PullRequestUpdate

	for hasNext := true; hasNext; hasNext = (next != "") {
		var us []PullRequestUpdate
		us, next = GetPullRequestActivity(*c, r.Name, pr.ID, next)

		for _, u := range us {
			// Break the loop once it is older than the last updated time
			if u.Date.Before(pr.LastUpdated) {
				next = ""
				break
			}

			// Ignore approvals if this user is only a reviewer
			if pr.Author.UUID != c.UserUUID && u.ActivityType == Approval {
				continue
			}

			latest = append(latest, u)
		}
	}

	pr.Updates = append(latest, pr.Updates...)
	pr.LastUpdated = time.Now()
	pr.UnreadUpdatesCount = pr.CountUnread(c)
}
