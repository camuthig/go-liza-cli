package cmd

import (
	"github.com/spf13/cobra"
)

// unwatchCmd represents the unwatch command
var unwatchCmd = &cobra.Command{
	Use:   "unwatch [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Remove a repository from your watched list.",
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()

		n := args[0]

		delete(c.Repositories, n)

		c.Write()
	},
}

func init() {
	rootCmd.AddCommand(unwatchCmd)
}
