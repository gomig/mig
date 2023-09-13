package main

import (
	"github.com/gomig/cli"
	"github.com/gomig/mig/commands"
)

func main() {
	cli := cli.NewCLI("gomig", "GoMig framework cli tools")
	cli.AddCommand(commands.VersionCommand)
	cli.AddCommand(commands.NewCommand)
	cli.Run()
}
