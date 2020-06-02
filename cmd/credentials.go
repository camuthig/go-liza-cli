package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// credentialsCmd represents the credentials command
var credentialsCmd = &cobra.Command{
	Use:   "credentials [username] [token]",
	Args:  cobra.ExactArgs(2),
	Short: "Configure the credentials for your user.",
	Long: `Configure the credentials for your user.

Liza uses an app password to access BitBucket. You can create the password here:

https://bitbucket.org/account/settings/app-passwords/

Liza needs access to:
	- Account read
	- Repositories read
	- Pull requests read`,
	Run: func(cmd *cobra.Command, args []string) {
		var c Config
		if err := c.Parse(); err != nil {
			fmt.Println("Unable to read config")
		}

		fmt.Println(cfgFile)

		c.Username = args[0]
		c.Token = args[1]

		if err := c.Write(); err != nil {
			fmt.Println("Unable to save config")
		}
	},
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
}
