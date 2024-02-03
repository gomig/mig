package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// UpdateCommand update cli to last version
var UpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "update cli to last version",
	Run: func(c *cobra.Command, args []string) {
		cmd := exec.Command("go", "install", "github.com/gomig/mig")
		fmt.Printf("Install github.com/gomig/mig: ")
		if err := cmd.Run(); err != nil {
			fmt.Printf("FAILED!\n")
			fmt.Println(err)
		} else {
			fmt.Printf("OK\n")
		}
	},
}
