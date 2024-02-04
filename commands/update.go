package commands

import (
	"fmt"
	"os/exec"

	"github.com/gomig/mig/helpers"
	"github.com/spf13/cobra"
)

// UpdateCommand update cli to last version
var UpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "update cli to last version",
	Run: func(c *cobra.Command, args []string) {
		cmd := exec.Command("go", "install", "github.com/gomig/mig")
		fmt.Println("Install github.com/gomig/mig")
		if err := cmd.Run(); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
		} else {
			fmt.Println(helpers.SuccessF("DONE"))
		}
	},
}
