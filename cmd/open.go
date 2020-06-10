package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:       "open",
	ValidArgs: []string{"name", "id"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
			return err
		}

		if len(args) > 0 {
			matched, err := regexp.Match(`^\S+\/\S+$`, []byte(args[0]))
			if err != nil || !matched {
				return fmt.Errorf("unabled to parse %s as a repository name", args[0])
			}
		}

		if len(args) == 2 {
			if _, err := strconv.Atoi(args[1]); err != nil {
				return fmt.Errorf("unabled to parse %s as pull request ID", args[1])
			}
		}

		return nil
	},
	Short: "Open a pull request in your browser and mark it read",
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()

		if len(args) == 2 {
			name := args[0]
			id, err := strconv.Atoi(args[1])

			if err != nil {
				fmt.Printf("Invalid ID %s", args[1])
				os.Exit(1)
			}
			openPullRequest(&c, name, id)
		}

		var rName *string
		if len(args) > 0 {
			rName = &args[0]
		}

		if len(args) > 1 {
			prID, _ := strconv.Atoi(args[1])

			openPullRequest(&c, *rName, prID)

			return
		}

		promptPullRequest(&c, rName)

	},
}

func promptPullRequest(c *Config, rName *string) {
	prs := c.AllPullRequests(rName)
	searcher := func(input string, index int) bool {
		pr := prs[index]
		return strings.Contains(strings.ToLower(pr.Title), strings.ToLower(input))
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Repository.Name }}: {{ .Title }}",
		Inactive: "  {{ .Repository.Name }}: {{ .Title }}",
		Selected: fmt.Sprintf(`%s {{ .Repository.Name }}: {{ .Title }}`, promptui.IconGood),
		Details: `Repository: {{.Repository.Name}}
ID: {{.ID}}
Link: {{.Links.HTML.Href}}
Title: {{.Title}}`,
	}
	prompt := promptui.Select{
		Label:     "Pull Requests",
		Size:      1,
		Templates: templates,
		Searcher:  searcher,
		IsVimMode: true,
		Items:     prs,
	}

	i, _, err := prompt.Run()

	if err != nil && err.Error() == "^C" {
		// TODO Is this really the best way to handle signals? Probably not.
		return
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	openPullRequest(c, prs[i].Repository.Name, prs[i].ID)

	c.Write()
}

func openPullRequest(c *Config, rName string, id int) {
	r, found := c.Repositories[rName]

	if !found {
		fmt.Printf("Not watching repository %s", rName)
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
}

func init() {
	rootCmd.AddCommand(openCmd)
}
