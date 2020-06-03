package cmd

import (
	"fmt"
	"os"

	"github.com/ktrysmt/go-bitbucket"
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
		c := ParseConfig()

		c.Username = args[0]
		c.Token = args[1]

		b := bitbucket.NewBasicAuth(c.Username, c.Token)

		p, err := b.User.Profile()

		if err != nil {
			fmt.Println("Unable to log into BitBucket")
			fmt.Println(err)
			os.Exit(1)
		}

		j := p.(map[string]interface{})

		c.UserUUID = j["uuid"].(string)

		c.Write()
	},
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
}
