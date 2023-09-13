package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/gomig/mig/internal/helpers"
	"github.com/gomig/utils"
	"github.com/spf13/cobra"
)

// NewCommand create new app
var NewCommand = new(cobra.Command)

func init() {
	NewCommand.Use = "new [app name]"
	NewCommand.Short = "create new gomig app"
	NewCommand.Args = cobra.MinimumNArgs(1)
	NewCommand.Run = func(cmd *cobra.Command, args []string) {
		// name
		name := args[0]
		root := "./" + name

		if exists, _ := utils.FileExists(root); exists {
			helpers.Handle(fmt.Sprintf("directory %s exists!\n", name))
		}

		// download and extract template
		func(basePath string) {
			dest := basePath + "/gomig.zip"
			helpers.Handle(helpers.Download("https://github.com/gomig/boilerplate/archive/latest.zip", dest))
			helpers.Handle(helpers.Unzip(dest, path.Join("./", name)))
			helpers.Handle(os.Remove(dest))
		}(root)

		// Run wizard
		w := runWizard()
		setup(name, w)
		fmt.Printf("\n\n")

		// Clean template files
		func(basePath string) {
			files := utils.FindFile(basePath, ".tpl.*")
			for _, f := range files {
				os.Remove("./" + f)
			}
		}(root)

		// Tidy app
		func(basePath string) {
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = basePath
			fmt.Printf("tidy app: ")
			if err := cmd.Run(); err != nil {
				fmt.Printf("FAILED!\n")
				fmt.Println(err)
			} else {
				fmt.Printf("OK!\n")
			}
		}(root)

		// Format app
		func(basePath string) {
			cmd := exec.Command("go", "fmt", "./...")
			cmd.Dir = basePath
			fmt.Printf("formatting: ")
			if err := cmd.Run(); err != nil {
				fmt.Printf("FAILED!\n")
				fmt.Println(err)
			} else {
				fmt.Printf("OK!\n")
			}
		}(root)

		if w.Result("git") == "y" {
			// init git
			func(basePath string) {
				cmd := exec.Command("git", "init")
				cmd.Dir = basePath
				fmt.Printf("init git: ")
				if err := cmd.Run(); err != nil {
					fmt.Printf("FAILED!\n")
					fmt.Println(err)
				} else {
					fmt.Printf("OK!\n")
				}
			}(root)
			// commit init
			func(basePath string) {
				addCmd := exec.Command("git", "add", ".")
				addCmd.Dir = basePath
				fmt.Printf("commit: ")
				if err := addCmd.Run(); err != nil {
					fmt.Printf("FAILED!\n")
					fmt.Println(err)
				} else {
					cmtCmd := exec.Command("git", "commit", "-m", "init")
					cmtCmd.Dir = basePath
					if err := cmtCmd.Run(); err != nil {
						fmt.Printf("FAILED!\n")
						fmt.Println(err)
					} else {
						fmt.Printf("OK!\n")
					}
				}
			}(root)
		}

		// Final message
		fmt.Printf("\nApp created.\nEnjoy it!\n\n")
	}
}
