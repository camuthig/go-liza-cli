package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Add a repository to your watched list.",
	Run: func(cmd *cobra.Command, args []string) {
		n := args[0]
		parts := strings.SplitN(args[0], "/", 2)
		c := ParseConfig()

		if _, ok := c.Repositories[n]; ok {
			fmt.Printf("Already watching %s", n)
			return
		}

		b := bitbucket.NewBasicAuth(c.Username, c.Token)

		opts := &bitbucket.RepositoryOptions{
			Owner:    parts[0],
			RepoSlug: parts[1],
		}

		r, err := b.Repositories.Repository.Get(opts)

		if err != nil {
			fmt.Printf("Unable to find repository %s\n", n)
			fmt.Println(err)
			os.Exit(1)
		}

		if c.Repositories == nil {
			c.Repositories = make(map[string]*Repository)
		}

		c.Repositories[r.Full_name] = &Repository{
			Name: r.Full_name,
		}

		c.Write()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
