package commands

import (
	"fmt"
	"io/fs"
	"os/exec"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gomig/mig/app"
	"github.com/gomig/mig/helpers"
	"github.com/spf13/cobra"
)

// NewCommand create new app
var NewCommand = new(cobra.Command)

func init() {
	NewCommand.Use = "new [app name]"
	NewCommand.Short = "create new template based app from github repository"
	NewCommand.Args = cobra.MinimumNArgs(1)
	NewCommand.Flags().StringP("branch", "b", "master", "branch name")
	NewCommand.Flags().StringP("auth", "a", "", "auth information")
	NewCommand.Run = func(cmd *cobra.Command, args []string) {
		// init
		cli := new(app.Mig)
		auth := new(app.Authentications)
		cli.Init()
		auth.Init()
		auth.Read()

		// get repository
		cli.AddQuestion(app.Question{
			Name:        "repo",
			Description: "Enter your github repository",
			Default:     "gomig/boilerplate",
		})
		cli.Start()

		// Parse parameters
		repo := helpers.NormalizeRepo(cli.Result("repo"))
		name := strings.TrimSpace(args[0])
		branch := cmd.Flag("branch").Value.String()
		var cred *http.BasicAuth
		if key := cmd.Flag("auth").Value.String(); key != "" {
			if c := auth.Find(key); c == nil {
				fmt.Printf("Failed: %s auth not registered!\n", key)
				return
			} else {
				cred = &http.BasicAuth{Username: c.Username, Password: c.Token}
			}
		}

		// Fetch repository
		fmt.Println("Fetching Repository...")
		source := memfs.New()
		helpers.HandleV(git.Clone(memory.NewStorage(), source, &git.CloneOptions{
			Auth:          cred,
			URL:           repo,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName(branch),
		}))
		fmt.Println("")

		// Parse app configuration
		cli.Init()
		if config, err := util.ReadFile(source, "/mig.json"); err == nil {
			helpers.Handle(cli.Parse(config))
		} else if err != fs.ErrNotExist {
			helpers.Handle(err)
		}
		cli.AddIgnore("mig.json")
		if name := cli.Name(); name != "" {
			fmt.Println(name)
		}
		if intro := cli.Intro(); intro != "" {
			fmt.Println(intro)
		}
		cli.Start()

		// Proccess files
		util.Walk(source, source.Root(), func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !cli.ShouldIgnore(path) {
				content := helpers.HandleV(util.ReadFile(source, path))
				helpers.Handle(cli.Compile(helpers.NormalizePath(path), content))
			}
			return nil
		})

		// Create out filesystem
		distFS := osfs.New(name, osfs.WithBoundOS())
		for file, content := range cli.Compiled() {
			f := helpers.HandleV(distFS.Create(file))
			helpers.HandleV(f.Write([]byte(content)))
			f.Close()
		}

		// Run post scripts
		base := "./" + name
		for _, script := range cli.Scripts() {
			command := ""
			args := make([]string, 0)
			if parts := strings.Split(script, " "); len(parts) == 0 {
				continue
			} else {
				command = parts[0]
				if len(parts) > 1 {
					args = append(args, parts[1:]...)
				}
			}
			cmd := exec.Command(command, args...)
			cmd.Dir = base
			fmt.Printf("Run (%s): ", script)
			if err := cmd.Run(); err != nil {
				fmt.Printf("FAILED!\n")
				fmt.Println(err)
			} else {
				fmt.Printf("OK\n")
			}
		}

		fmt.Println("App Created")
	}
}
