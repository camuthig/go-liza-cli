package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// watchedCmd represents the watched command
var watchedCmd = &cobra.Command{
	Use:   "watched",
	Short: "List the repositories you are currently watching.",
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()

		for _, r := range c.Repositories {
			fmt.Println(r.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
