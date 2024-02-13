package commands

import (
	"fmt"

	"github.com/gomig/mig/helpers"
	"github.com/spf13/cobra"
)

// VersionCommand get cli version
var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "get mig cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s\n", helpers.SuccessF("2.2.0"))
	},
}
