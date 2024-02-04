package commands

import (
	"fmt"

	"github.com/gomig/mig/app"
	"github.com/gomig/mig/helpers"
	"github.com/spf13/cobra"
)

// UnAuthCommand delete auth from app
var UnAuthCommand = &cobra.Command{
	Use:   "unauth [key]",
	Short: "delete auth info",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		auth := new(app.Authentications)
		auth.Init()
		if err := auth.Read(); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		auth.Delete(args[0])
		if err := auth.Write(); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		fmt.Println(helpers.SuccessF("Auth Deleted"))
	},
}
