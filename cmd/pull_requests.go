package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var ForAllPullRequests bool

func ValidatePullRequestArgs() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
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
	}
}

func RunForPullRequests(action func(c *Config, pr *PullRequestWithRepository)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		c := ParseConfig()

		if len(args) == 2 {
			name := args[0]
			id, err := strconv.Atoi(args[1])

			if err != nil {
				fmt.Printf("Invalid ID %s", args[1])
				os.Exit(1)
			}
			action(&c, getPullRequest(&c, name, id))
		}

		var rName *string
		if len(args) > 0 {
			rName = &args[0]
		}

		if len(args) > 1 {
			prID, _ := strconv.Atoi(args[1])

			action(&c, getPullRequest(&c, *rName, prID))

			c.Write()
			return
		}

		if ForAllPullRequests {
			for _, pr := range c.AllPullRequests(rName) {
				action(&c, &pr)
			}

			c.Write()
			return
		}

		pr := promptPullRequest(&c, rName)

		if pr == nil {
			return
		}

		action(&c, pr)

		c.Write()
	}
}

func getUnreadPrompt(pr PullRequestWithRepository) string {
	s := fmt.Sprintf("[%d]", pr.UnreadUpdatesCount)
	if pr.UnreadUpdatesCount > 0 {
		return promptui.Styler(promptui.FGGreen)(s)
	}

	return s
}

func selectPullRequests(c *Config, repo *string) (*PullRequestWithRepository, bool) {
	prs := c.AllPullRequests(repo)

	searcher := func(input string, index int) bool {
		pr := prs[index]
		return strings.Contains(strings.ToLower(pr.Title), strings.ToLower(input))
	}

	funcs := promptui.FuncMap
	funcs["getUnreadPrompt"] = getUnreadPrompt

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Repository.Name }}{{ getUnreadPrompt . }}: {{ .Title }}",
		Inactive: "  {{ .Repository.Name }}{{ getUnreadPrompt .}}: {{ .Title }}",
		Selected: fmt.Sprintf(`%s {{ .Repository.Name }}: {{ .Title }}`, promptui.IconGood),
		Details: `Repository: {{.Repository.Name}}
ID: {{.ID}}
Link: {{.Links.HTML.Href}}
Title: {{.Title}}`,
		FuncMap: funcs,
	}

	prompt := promptui.Select{
		Label:     "Pull Requests",
		Size:      10,
		Templates: templates,
		Searcher:  searcher,
		IsVimMode: true,
		Items:     prs,
	}

	i, _, err := prompt.Run()

	if err != nil && err.Error() == "^C" {
		// TODO Is this really the best way to handle signals? Probably not.
		return nil, false
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, false
	}

	return &prs[i], true
}

func promptPullRequest(c *Config, rName *string) *PullRequestWithRepository {
	pr, selected := selectPullRequests(c, rName)

	if !selected {
		return nil
	}

	return pr
}

func getPullRequest(c *Config, rName string, id int) *PullRequestWithRepository {
	repo := getRepo(c, rName)

	pr, found := repo.PullRequests[id]

	if !found {
		fmt.Printf("Pull request %d not found", id)
		os.Exit(1)
	}

	return &PullRequestWithRepository{PullRequest: pr, Repository: repo}
}

func getRepo(c *Config, rName string) *Repository {
	r, found := c.Repositories[rName]

	if !found {
		fmt.Printf("Not watching repository %s", rName)
		os.Exit(1)
	}

	return r
}
