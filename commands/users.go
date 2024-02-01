package commands

import (
	"fmt"

	"github.com/gomig/mig/app"
	"github.com/spf13/cobra"
)

// UsersCommand show users list
var UsersCommand = &cobra.Command{
	Use:   "users",
	Short: "show users list",
	Run: func(cmd *cobra.Command, args []string) {
		auth := new(app.Authentications)
		auth.Init()
		if err := auth.Read(); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Key                           User")
		fmt.Println("+----------------------------+----")
		for _, cred := range auth.Credentials() {
			fmt.Printf("%-30s%s\n", cred.Key, cred.Username)
		}
	},
}
