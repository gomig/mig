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
	"github.com/gomig/crypto"
	"github.com/gomig/mig/app"
	"github.com/gomig/mig/helpers"
	"github.com/google/uuid"
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

		// Validate params
		if cred := auth.Find(cmd.Flag("auth").Value.String()); cred == nil && cmd.Flag("auth").Value.String() != "" {
			fmt.Println(helpers.ErrorF("auth not registered!"))
			return
		}

		// get repository
		cli.AddRule(app.Rule{
			Name:        "repo",
			Description: "Enter template repository",
			Default:     "gomig/boilerplate",
		})
		cli.Start()

		// Parse parameters
		repo := helpers.NormalizeRepo(cli.Result("repo"))
		name := strings.TrimSpace(args[0])
		branch := cmd.Flag("branch").Value.String()
		var cred *http.BasicAuth
		if c := auth.Find(cmd.Flag("auth").Value.String()); c != nil {
			cred = &http.BasicAuth{Username: c.Username, Password: c.Token}
		}

		// Fetch repository
		fmt.Println(helpers.Format("Fetching Repository...", helpers.ITALIC, helpers.BLUE))
		source := memfs.New()
		if _, err := git.Clone(memory.NewStorage(), source, &git.CloneOptions{
			Auth:          cred,
			URL:           repo,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName(branch),
		}); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}
		fmt.Println("")

		// Parse app configuration
		cli.Init()
		if config, err := util.ReadFile(source, "/mig.json"); err == nil {
			if err := cli.Parse(config); err != nil {
				fmt.Println(helpers.ErrorF("Can't parse config file!"))
				return
			}
		} else if err != fs.ErrNotExist {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}
		cli.AddIgnore("mig.json")
		if name := cli.Name(); name != "" {
			fmt.Println(helpers.Format(name, helpers.BOLD, helpers.BLUE))
		}
		if intro := cli.Intro(); intro != "" {
			fmt.Println(intro)
		}
		cli.Start()

		// Add global values
		cli.AddRule(app.Rule{
			Name:        "name",
			Default:     name,
			Placeholder: "__name__",
			Description: "Application name",
		})
		cli.AddRule(app.Rule{
			Name:        "key",
			Placeholder: "__key__",
			Description: "Application key",
		})
		cli.AddRule(app.Rule{
			Name:        "token",
			Placeholder: "__token__",
			Description: "Application access token",
		})
		cli.AddRule(app.Rule{
			Name:        "uuid",
			Placeholder: "__uuid__",
			Description: "uuid",
		})

		// Generate keys and fill global variable
		cli.AddAnswer("name", name)
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.SHA3256); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		} else {
			cli.AddAnswer("key", key)
		}
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.SHA3256); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		} else {
			cli.AddAnswer("token", key)
		}
		cli.AddAnswer("uuid", uuid.New().String())

		// Proccess files
		if err := util.Walk(source, source.Root(), func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !cli.ShouldIgnore(path) && err == nil {
				if content, err := util.ReadFile(source, path); err != nil {
					return err
				} else {
					if err := cli.Compile(path, content); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			fmt.Println(helpers.ErrorF(err.Error()))
			return
		}

		// Create out filesystem
		distFS := osfs.New(name, osfs.WithBoundOS())
		for file, content := range cli.Compiled() {
			if f, err := distFS.Create(file); err != nil {
				fmt.Println(helpers.ErrorF(err.Error()))
				return
			} else {
				if _, err := f.Write([]byte(content)); err != nil {
					f.Close()
					fmt.Println(helpers.ErrorF(err.Error()))
					return
				}
				f.Close()
			}
		}

		// Run post scripts
		base := "./" + name
		for _, parts := range cli.Scripts() {
			command := ""
			args := make([]string, 0)
			if len(parts) == 0 {
				continue
			} else {
				command = parts[0]
				if len(parts) > 1 {
					args = append(args, parts[1:]...)
				}
			}
			cmd := exec.Command(command, args...)
			cmd.Dir = base
			fmt.Printf("Run (%s)\n", helpers.Format(strings.Join(parts, " "), helpers.ITALIC, helpers.BLUE))
			if err := cmd.Run(); err != nil {
				fmt.Println(helpers.ErrorF(err.Error()))
			} else {
				fmt.Println(helpers.SuccessF("DONE"))
			}
		}

		fmt.Println("")
		fmt.Println(helpers.Format("Project Created", helpers.BOLD, helpers.GREEN))
		if message := cli.Message(); message != "" {
			fmt.Println(message)
		}
	}
}
