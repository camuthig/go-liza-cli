/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Use:   "watch",
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
			c.Repositories = make(map[string]Repository)
		}

		c.Repositories[r.Full_name] = Repository{
			Name: r.Full_name,
		}

		c.Write()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
