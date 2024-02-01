package main

import (
	"github.com/gomig/cli"
	"github.com/gomig/mig/commands"
)

func main() {
	cli := cli.NewCLI("mig", "create new template based app")
	cli.AddCommand(commands.VersionCommand)
	cli.AddCommand(commands.AuthCommand)
	cli.AddCommand(commands.UnAuthCommand)
	cli.AddCommand(commands.UsersCommand)
	cli.AddCommand(commands.NewCommand)
	cli.Run()
}
