package commands

import "github.com/gomig/mig/internal/questions"

func runWizard() *questions.Wizard {
	// Run wizard
	wiz := new(questions.Wizard)
	wiz.Init()

	wiz.AddQuestion(questions.Question{
		Name:        "description",
		Description: "description",
		Default:     "",
		Valids:      []string{},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "namespace",
		Description: "namespace",
		Default:     "github.com/anonymous",
		Valids:      []string{},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "locale",
		Description: "enter app locale",
		Default:     "fa",
		Valids:      []string{},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "config",
		Description: "select app configuration manager (env/json/memory)",
		Default:     "json",
		Valids:      []string{"env", "json", "memory"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "cache",
		Description: "select cache manager (file/redis)",
		Default:     "redis",
		Valids:      []string{"file", "redis"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "translator",
		Description: "select translator driver (json/memory)",
		Default:     "memory",
		Valids:      []string{"json", "memory"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "database",
		Description: "database driver (mongo) (mongo/mysql/none)",
		Default:     "mongo",
		Valids:      []string{"mongo", "mysql", "none"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "web",
		Description: "include web (gofiber) (y/n)",
		Default:     "y",
		Valids:      []string{"y", "n"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "template",
		Description: "include template engine (html/template) (y/n)",
		Default:     "y",
		Valids:      []string{"y", "n"},
	})

	wiz.AddQuestion(questions.Question{
		Name:        "git",
		Description: "init git repository (y/n)",
		Default:     "y",
		Valids:      []string{"y", "n"},
	})

	wiz.Start()
	return wiz
}
