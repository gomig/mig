package commands

import (
	"fmt"

	"github.com/gomig/mig/app"
	"github.com/gomig/mig/helpers"
	"github.com/spf13/cobra"
)

// AuthCommand add auth to app
var AuthCommand = &cobra.Command{
	Use:   "auth [key] [user] [access token]",
	Short: "add auth info",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		auth := new(app.Authentications)
		auth.Init()
		if err := auth.Read(); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		auth.Add(args[0], args[1], args[2])
		if err := auth.Write(); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		fmt.Println(helpers.SuccessF("Auth Added"))
	},
}
