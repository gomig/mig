package commands

import (
	"fmt"
	"os"

	"github.com/gomig/mig/app"
	"github.com/gomig/mig/helpers"
	"github.com/jedib0t/go-pretty/v6/table"
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
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Authenticate Key", "Username"})
		for i, cred := range auth.Credentials() {
			t.AppendRow([]interface{}{i + 1, cred.Key, cred.Username})
		}
		t.SetStyle(table.StyleLight)
		t.Render()
	},
}
